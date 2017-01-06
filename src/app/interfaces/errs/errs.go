package errs

import (
	"fmt"

	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Error Codes
const (
	NotImplementedError uint16 = iota
	InternalServerError
	TokenExpired
)

// todo: remove it
var (
	ErrTokenExpired = New(TokenExpired, 401, "Token expired")
)

func New(code uint16, hCode int, msg string, args ...interface{}) *Error {
	return &Error{Msg: msg, Code: code, HTTPCode: hCode, Args: args}
}

func Unauthorized(msg string, args ...interface{}) *Error {
	return New(NotImplementedError, http.StatusUnauthorized, msg, args)
}

func BadRequest(msg string, args ...interface{}) *Error {
	return New(NotImplementedError, http.StatusBadRequest, msg, args)
}

// IsTokenExpiredErr checks given error is jwt expired token error
func IsTokenExpiredErr(err error) bool {
	vErr, ok := Cause(err).(*jwt.ValidationError)
	if ok && vErr.Errors == jwt.ValidationErrorExpired {
		return true
	}
	return false
}

func IsTokenValidationErr(err error) bool {
	_, ok := Cause(err).(*jwt.ValidationError)
	return ok
}

// Error is application error
type Error struct {
	Inner    error
	Msg      string
	Code     uint16
	HTTPCode int
	Args     []interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf(e.Msg, e.Args...)
}

func (e Error) SetInner(err error) *Error {
	e.Inner = err
	return &e
}

var NewWithStack = errors.Errorf
var Wrap = errors.WithStack
var WrapMsg = errors.Wrapf
var Cause = errors.Cause
