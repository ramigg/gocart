package usecases

import (
	"app"
	"app/infra/social"
	"app/infra/social/providers"
	"app/interfaces/errs"
)

func NewSocialAuth() *socialAuth {
	return &socialAuth{}
}

type socialAuth struct{}

func (sa *socialAuth) GetUserFromFacebook(token string) (*app.User, error) {
	fb := social.New(providers.NewFacebook(token))
	su, err := fb.GetUser()
	if err != nil {
		if err, ok := err.(*social.ApiErr); ok {
			return nil, errs.BadRequest("user can't get from facebook").SetInner(err)
		}
		return nil, err
	}

	var u app.User
	u.Email = su.Email
	u.FirstName = su.FirstName
	u.LastName = su.LastName
	return &u, nil
}

func NewMockSocialAuth() *mockSocialAuth {
	return &mockSocialAuth{}
}

type mockSocialAuth struct{}

func (msa *mockSocialAuth) GetUserFromFacebook(token string) (*app.User, error) {
	validToken := "valid token"
	if token != validToken {
		return nil, errs.BadRequest("api error")
	}

	var u app.User
	u.Email = "user@facebook.com"
	u.FirstName = "Mark"
	u.LastName = "Zuckerberg"
	return &u, nil
}
