package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/login", loginHandler)
	http.HandleFunc("/chat", chatPage)

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func chatPage(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		username = "Guest"
	}

	tmpl := template.Must(template.New("chat").Parse(`
		<!DOCTYPE html>
		<html>
		<head><title>Chat</title></head>
		<body>
			<h1>Welcome, {{.}}!</h1>
			<p>You are now in the chat room.</p>
		</body>
		</html>`))
	tmpl.Execute(w, username)
}
