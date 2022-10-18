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
}


type Sample struct {
	Data []ClientUser
	Title string
	AdminName string
	ErrMsg string
	Id string
	Email string
}
