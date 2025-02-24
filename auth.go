package main

import (
	"net/http"

	"github.com/stretchr/objx"
)

type chatUser struct {
	nickname string
	uniqueID string
	avatar   string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     "myuserid",
			"name":       username,
			"avatar_url": "myavatar",
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		w.Header().Set("Location", "/chat?user="+username)
		http.Redirect(w, r, "/chat?user="+username, http.StatusSeeOther)
		return
	}
	http.Error(w, "Invalid request", http.StatusBadRequest)
}
