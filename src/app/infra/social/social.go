package social

func New(p Provider) Provider {
	return p
}

type Provider interface {
	GetUser() (*User, error)
}

type User struct {
	Email     string
	FirstName string
	LastName  string
}

func NewApiErr(statusCode int, body []byte) *ApiErr {
	return &ApiErr{statusCode, body}
}

type ApiErr struct {
	statusCode int
	body       []byte
}

// todo: need better error msg
func (ae *ApiErr) Error() string {
	return "api error"
}
func (ae *ApiErr) StatusCode() int {
	return ae.statusCode
}
