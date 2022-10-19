package utils

import "github.com/gorilla/sessions"

var (
	UserKey   = []byte("for-user")
	UserStore = sessions.NewCookieStore(UserKey)
)
