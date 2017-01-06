package handlers

import (
	"app"
	"app/interfaces/errs"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"

	gores "gopkg.in/alioygur/gores.v1"
)

var (
	errEmailExists    = errs.BadRequest("email address already exists")
	errEmailNotExists = errs.BadRequest("email address not exists")
	errWrongCred      = errs.Unauthorized("wrong credentials")
	errInactiveUser   = errs.Unauthorized("inactive user")
	errInvalidToken   = errs.BadRequest("invalid token")
)

type userRepo interface {
	OneByEmail(string) (*app.User, error)
	ExistsByEmail(string) (bool, error)
	UpdateFields(map[string]interface{}) error
	Create(*app.User) error
	app.DBNotFoundErrChecker
}

type socialAuth interface {
	GetUserFromFacebook(string) (*app.User, error)
}

// NewAuthHandler instances new auth handler struct
func NewAuthHandler(ur userRepo, sa socialAuth) *authHandler {
	return &authHandler{ur, sa}
}

// AuthHandler struct
type authHandler struct {
	ur userRepo
	sa socialAuth
	// sendPasswordResetLink func(m app.MailSender, to, link string) error
}

// SetRoutes sets this module's routes
func (ah *authHandler) SetRoutes(r *mux.Router) {
	r.Handle("/v1/auth/login", appHandler(ah.login)).Methods("POST")
	r.Handle("/v1/auth/register", appHandler(ah.register)).Methods("POST")
	r.Handle("/v1/auth/register-fb", appHandler(ah.registerFacebook)).Methods("POST")

	r.Handle("/v1/password/forgot", appHandler(ah.forgotPassword)).Methods("POST")
	r.Handle("/v1/password/reset", appHandler(ah.resetPassword)).Methods("POST")
}

func (ah *authHandler) login(w http.ResponseWriter, r *http.Request) error {
	f := new(loginForm)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	u, err := ah.ur.OneByEmail(f.Email)
	if err != nil {
		if ah.ur.IsNotFoundErr(err) {
			return errWrongCred
		}
		return err
	}

	if !u.IsCredentialsVerified(f.Password) {
		return errWrongCred
	}

	if !u.IsActivated {
		return errInactiveUser
	}

	token, err := u.CreateJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusOK, tokenRes{token})
}

func (ah *authHandler) register(w http.ResponseWriter, r *http.Request) error {
	f := new(registerForm)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	if err := errs.CheckEmail(f.Email); err != nil {
		return err
	}
	if err := errs.CheckPassword(f.Password); err != nil {
		return err
	}

	// check for email
	exists, err := ah.ur.ExistsByEmail(f.Email)
	if err != nil {
		return err
	} else if exists {
		return errEmailExists
	}

	var usr app.User
	usr.FirstName = f.FirstName
	usr.LastName = f.LastName
	usr.Email = f.Email
	usr.IsActivated = true
	usr.SetPassword(f.Password)

	if err := ah.ur.Create(&usr); err != nil {
		return err
	}

	token, err := usr.CreateJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusCreated, tokenRes{token})
}

func (ah *authHandler) forgotPassword(w http.ResponseWriter, r *http.Request) error {
	f := new(forgotPasswordForm)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	if err := errs.CheckEmail(f.Email); err != nil {
		return err
	}

	if f.Link == "" {
		f.Link = os.Getenv("PASSWORD_RESET_URL")
	}
	resetURL, err := url.Parse(f.Link)
	if err != nil {
		return errs.BadRequest("invalid url").SetInner(err)
	}

	u, err := ah.ur.OneByEmail(f.Email)
	if err != nil {
		if ah.ur.IsNotFoundErr(err) {
			return errEmailNotExists
		}
		return err
	}

	tokenString, err := u.GenResetPasswordToken()
	if err != nil {
		return err
	}

	q := resetURL.Query()
	q.Set("token", tokenString)
	resetURL.RawQuery = q.Encode()

	body := fmt.Sprintf("Please click below link to reset your password <br/> %s", resetURL.String())

	// todo: Send email here
	_ = body

	gores.NoContent(w)
	return nil
}

func (ah *authHandler) resetPassword(w http.ResponseWriter, r *http.Request) error {
	f := new(resetPasswordForm)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	if err := errs.CheckEmail(f.Email); err != nil {
		return err
	}
	if err := errs.CheckPassword(f.Password); err != nil {
		return err
	}

	u, err := ah.ur.OneByEmail(f.Email)
	if err != nil {
		if ah.ur.IsNotFoundErr(err) {
			return errEmailNotExists
		}
		return err
	}

	if err := u.IsResetPasswordTokenValid(f.Token); err != nil {
		if errs.IsTokenValidationErr(err) {
			return errInvalidToken.SetInner(err)
		}
		return err
	}

	u.SetPassword(f.Password)
	if err := ah.ur.UpdateFields(map[string]interface{}{"Password": u.Password}); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}

func (ah *authHandler) registerFacebook(w http.ResponseWriter, r *http.Request) error {
	f := new(registerFacebook)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	u, err := ah.sa.GetUserFromFacebook(f.AccessToken)
	if err != nil {
		return err
	}

	existsUser, err := ah.ur.OneByEmail(u.Email)
	if err == nil {
		u = existsUser
	} else if ah.ur.IsNotFoundErr(err) {
		u.IsActivated = true
		if err := ah.ur.Create(u); err != nil {
			return err
		}
	} else {
		return err
	}

	token, err := u.CreateJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusCreated, tokenRes{token})
}

type tokenRes struct {
	Token string `json:"token"`
}

type updateUserForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type registerForm struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsActivated bool   `json:"isActivated"`
}

type registerFacebook struct {
	AccessToken string `json:"accessToken"`
}

type forgotPasswordForm struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

type resetPasswordForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
