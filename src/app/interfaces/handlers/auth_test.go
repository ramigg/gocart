package handlers

import (
	"app/interfaces/repos/mockdb"
	"net/http"
	"testing"

	"app"

	"encoding/json"

	"app/usecases"

	"github.com/gorilla/mux"
)

func TestAuthHandler_login(t *testing.T) {
	h := mux.NewRouter()

	// new user repo
	ur := &mockdb.User{}

	// add a valid user.
	var u app.User
	u.Email = "activeuser@gmail.com"
	u.SetPassword("good password")
	u.IsActivated = true
	if err := ur.Create(&u); err != nil {
		t.Fatal(err)
	}
	// add a inactive user.
	var u2 app.User
	u2.Email = "inactiveuser@gmail.com"
	u2.SetPassword("good password")
	u2.IsActivated = false
	if err := ur.Create(&u2); err != nil {
		t.Fatal(err)
	}

	ah := NewAuthHandler(ur, usecases.NewMockSocialAuth())
	ah.SetRoutes(h)

	var (
		goodLoginCred, _    = json.Marshal(loginForm{"activeuser@gmail.com", "good password"})
		badLoginCred, _     = json.Marshal(loginForm{"wrong@mail.address", "wrong password"})
		inactiveUserCred, _ = json.Marshal(loginForm{"inactiveuser@gmail.com", "good password"})
	)

	testCases := []testCase{
		{"login with no cred.", "/v1/auth/login", "POST", nil, http.StatusUnauthorized, nil},
		{"login with bad cred.", "/v1/auth/login", "POST", badLoginCred, http.StatusUnauthorized, nil},
		{"login with inactive user", "/v1/auth/login", "POST", inactiveUserCred, http.StatusUnauthorized, nil},
		{"login with good cred.", "/v1/auth/login", "POST", goodLoginCred, http.StatusOK, nil},
	}

	runHandlerTestCases(testCases, h, t)
}

func TestAuthHandler_register(t *testing.T) {
	h := mux.NewRouter()

	// new user repo
	ur := &mockdb.User{}

	ah := NewAuthHandler(ur, usecases.NewMockSocialAuth())
	ah.SetRoutes(h)

	var (
		goodRegisterParams, _ = json.Marshal(registerForm{Email: "newuser@gmail.com", Password: "new password"})
	)

	testCases := []testCase{
		{"register with no params", "/v1/auth/register", "POST", nil, http.StatusBadRequest, nil},
		{"register with good params", "/v1/auth/register", "POST", goodRegisterParams, http.StatusCreated, nil},
		{"register with already existing email", "/v1/auth/register", "POST", goodRegisterParams, http.StatusBadRequest, nil},
	}

	runHandlerTestCases(testCases, h, t)
}

func TestAuthHandler_forgotPassword(t *testing.T) {
	h := mux.NewRouter()

	// new user repo
	ur := &mockdb.User{}

	// add a valid user.
	var u app.User
	u.Email = "activeuser@gmail.com"
	u.SetPassword("good password")
	u.IsActivated = true
	if err := ur.Create(&u); err != nil {
		t.Fatal(err)
	}

	ah := NewAuthHandler(ur, usecases.NewMockSocialAuth())
	ah.SetRoutes(h)

	var (
		goodLoginCred, _ = json.Marshal(loginForm{"activeuser@gmail.com", "goodpassword"})
		badLoginCred, _  = json.Marshal(loginForm{"wrong@mail.address", "wrong password"})
	)

	testCases := []testCase{
		{"forgot password with no params", "/v1/password/forgot", "POST", nil, http.StatusBadRequest, nil},
		{"forgot password with bad params", "/v1/password/forgot", "POST", badLoginCred, http.StatusBadRequest, nil},
		{"forgot password with good params", "/v1/password/forgot", "POST", goodLoginCred, http.StatusNoContent, nil},
	}

	runHandlerTestCases(testCases, h, t)
}

func TestAuthHandler_resetPassword(t *testing.T) {
	h := mux.NewRouter()

	// new user repo
	ur := &mockdb.User{}

	// add a valid user.
	var u app.User
	u.Email = "activeuser@gmail.com"
	u.SetPassword("good password")
	u.IsActivated = true
	if err := ur.Create(&u); err != nil {
		t.Fatal(err)
	}

	ah := NewAuthHandler(ur, usecases.NewMockSocialAuth())
	ah.SetRoutes(h)

	resetToken, err := u.GenResetPasswordToken()
	if err != nil {
		t.Fatalf("can't generate reset password token: %v", err)
	}

	var (
		goodParams, _   = json.Marshal(resetPasswordForm{"activeuser@gmail.com", "my new password", resetToken})
		invalidEmail, _ = json.Marshal(resetPasswordForm{"invalid@gmail.com", "my new password", resetToken})
		invalidToken, _ = json.Marshal(resetPasswordForm{"activeuser@gmail.com", "my new password", "bad token"})
	)

	testCases := []testCase{
		{"reset password with no params", "/v1/password/reset", "POST", nil, http.StatusBadRequest, nil},
		{"reset password with invalid email", "/v1/password/reset", "POST", invalidEmail, http.StatusBadRequest, nil},
		{"reset password with invalid token", "/v1/password/reset", "POST", invalidToken, http.StatusBadRequest, nil},
		{"reset password good params", "/v1/password/reset", "POST", goodParams, http.StatusNoContent, nil},
	}

	runHandlerTestCases(testCases, h, t)
}

func TestAuthHandler_registerFacebook(t *testing.T) {
	h := mux.NewRouter()

	// new user repo
	ur := &mockdb.User{}

	ah := NewAuthHandler(ur, usecases.NewMockSocialAuth())
	ah.SetRoutes(h)

	var (
		goodParams, _   = json.Marshal(registerFacebook{"valid token"})
		invalidToken, _ = json.Marshal(registerFacebook{"invalid token"})
	)

	testCases := []testCase{
		{"register facebook with no params", "/v1/auth/register-fb", "POST", nil, http.StatusBadRequest, nil},
		{"register facebook with invalid token", "/v1/auth/register-fb", "POST", invalidToken, http.StatusBadRequest, nil},
		{"register facebook good params", "/v1/auth/register-fb", "POST", goodParams, http.StatusCreated, nil},
	}

	runHandlerTestCases(testCases, h, t)
}
