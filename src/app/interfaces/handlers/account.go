package handlers

import (
	"app"
	"net/http"

	"app/interfaces/errs"

	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type accountRepo interface {
	ExistsByEmail(string) (bool, error)
	UpdateFields(map[string]interface{}) error
}

func NewAccount(ur accountRepo) *Account {
	return &Account{ur}
}

type Account struct {
	ur accountRepo
}

func (a *Account) SetRoutes(r *mux.Router, mid ...alice.Constructor) {
	h := alice.New(mid...)
	r.Handle("/v1/me", h.Then(appHandler(a.me))).Methods("GET")
	r.Handle("/v1/me", h.Then(appHandler(a.update))).Methods("PATCH")
}

func (a *Account) me(w http.ResponseWriter, r *http.Request) error {
	u := app.UserMustFromContext(r.Context())

	return gores.JSON(w, http.StatusOK, u)
}

func (a *Account) update(w http.ResponseWriter, r *http.Request) error {
	f := new(updateMeForm)
	if err := decodeReq(r, f); err != nil {
		return err
	}

	me := app.UserMustFromContext(r.Context())

	exists, err := a.ur.ExistsByEmail(f.Email)
	if err != nil {
		return err
	}

	if exists {
		return errEmailExists
	}

	fields := make(map[string]interface{})

	if f.Email != "" {
		if err := errs.CheckEmail(f.Email); err != nil {
			return err
		}
		fields["Email"] = f.Email
	}

	if f.Password != "" {
		if err := errs.CheckPassword(f.Password); err != nil {
			return err
		}
		me.SetPassword(f.Password)
		fields["Password"] = f.Password
	}

	if err := a.ur.UpdateFields(fields); err != nil {
		return err
	}

	gores.NoContent(w)
	return nil
}

type updateMeForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
