package mockdb

import "app"

type User struct {
	*Repo
	users []*app.User
}

func (ur *User) Create(u *app.User) error {
	ur.users = append(ur.users, u)
	return nil
}

func (ur *User) OneByEmail(email string) (*app.User, error) {
	for _, u := range ur.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errNotFound
}

func (ur *User) ExistsByEmail(email string) (bool, error) {
	for _, u := range ur.users {
		if u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (ur *User) UpdateFields(kv map[string]interface{}) error {
	return nil
}
