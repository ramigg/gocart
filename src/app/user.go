package app

import (
	"app/interfaces/errs"
	"context"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

var userContextKey contextKey = "user"

// User model
type User struct {
	Model       `storm:"inline"`
	FirstName   string `json:"firstName" fako:"first_name"`
	LastName    string `json:"lastName" fako:"last_name"`
	Email       string `json:"email" fako:"email_address" storm:"unique"`
	Password    string `json:"password" fako:"simple_password"`
	IsActivated bool   `json:"isActivated"`
	IsAdmin     bool   `json:"isAdmin"`
}

// SetPassword sets user's password
func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.Password = string(hashedPassword)
}

// IsCredentialsVerified matches given password with user's password
func (u *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

// CreateJWT creates a Javascript Web Token (JWT)
func (u *User) CreateJWT(secretKey string) (string, error) {
	// claims
	id := strconv.Itoa(int(u.ID))
	claims := jwt.MapClaims{
		"userID": id,
		"exp":    time.Now().Add(time.Hour * 6).Unix(),
	}

	// token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// the final token (hashed string)
	signedSecret, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errs.WrapMsg(err, "token can't signed")
	}

	return signedSecret, nil
}

func (u *User) GenResetPasswordToken() (string, error) {
	claims := jwt.MapClaims{"email": u.Email, "exp": time.Now().Add(time.Hour * 5).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.Password))
	if err != nil {
		return "", errs.WrapMsg(err, "token can't signed")
	}
	return tokenString, nil
}

func (u *User) IsResetPasswordTokenValid(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Password), nil
	})
	if err != nil {
		return err
	}

	email, ok := token.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return errs.NewWithStack("email can't get from token claims, token: %s", tokenStr)
	}

	if email != u.Email {
		return errs.NewWithStack("token's email and user's email aren't equal")
	}

	return nil
}

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext gets user from context
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// UserMustFromContext gets user from context. if can't make panic
func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}
