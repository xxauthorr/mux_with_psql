package utils

import "github.com/gorilla/sessions"

var (
	UserKey    = []byte("for-user")
	// AdminKey   = []byte("for-admin")
	UserStore  = sessions.NewCookieStore(UserKey)
	// AdminStore = sessions.NewCookieStore(AdminKey)
)
