package base

import "github.com/gorilla/sessions"

var SessionStore = sessions.NewCookieStore([]byte("session"))
