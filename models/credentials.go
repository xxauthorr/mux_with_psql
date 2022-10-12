package models

type Credentials struct {
	Errmsg string
	Header string
	Email  string
	LoggedIn bool
}

type ClientUser struct {
	Id       int32
	Email    string
	Hashpass string
}
