package providers

import (
	"app/infra/social"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewFacebook(token string) *facebook {
	return &facebook{token}
}

type facebook struct {
	accessToken string
}

func (fb *facebook) GetUser() (*social.User, error) {
	fields := "email,first_name,last_name"
	baseURL := fmt.Sprintf("https://graph.facebook.com/v2.7/me/?access_token=%s&fields=%s", fb.accessToken, fields)

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 199 {
		return nil, social.NewApiErr(resp.StatusCode, body)
	}

	return fb.mapUser(body)
}

func (fb *facebook) mapUser(body []byte) (*social.User, error) {
	var u social.User
	fbUser := struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{}
	if err := json.Unmarshal(body, fbUser); err != nil {
		return nil, err
	}

	u.Email = fbUser.Email
	u.FirstName = fbUser.Email
	u.LastName = fbUser.LastName
	return &u, nil

}
