package models

type Credentials struct {
	Errmsg   string
	Header   string
	Email    string
	LoggedIn bool
}

type ClientUser struct {
	Id       string
	Email    string
	Hashpass string
	ErrMsg   string
	Username string
	Title    string
}

type UserData []struct {
	id    string
	email string
}

type Sample struct {
	Data []ClientUser
	Title string
}
