package gormdb

import "app"

func NewUser(r *Repo) *User {
	return &User{r}
}

type User struct {
	*Repo
}

func (ur *User) OneByEmail(email string) (*app.User, error) {
	var u app.User
	return &u, ur.OneBy(&u, app.DBWhere{"Email": email})
}

func (ur *User) ExistsByEmail(email string) (bool, error) {
	var u app.User
	return ur.ExistsBy(&u, app.DBWhere{"Email": email})
}

func (ur *User) UpdateFields(kv map[string]interface{}) error {
	return ur.Repo.UpdateFields(&app.User{}, kv)
}

func (ur *User) Create(u *app.User) error {
	return ur.Store(u)
}
