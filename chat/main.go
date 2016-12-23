package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"work/goweb/chat/trace"
	"os"
)

func main() {
	addr := flag.String("addr", ":8080", "Apprication's Address")
	flag.Parse()

	//Gomniauth Set-up
	gomniauth.SetSecurityKey("dGhKIuhkmhtqxNNbUuiJNVtgun[;?<,miHNbgojj")
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		github.New(
			"aec2bc42117b12eb1328",
			"56282cd7e64284f019554221fad29d12aca52714",
			"http://localhost:8080/auth/callback/github"),
		google.New(
			"176223784373-bss719fm64msnfuqau908ql7l1es0dcs.apps.googleusercontent.com",
			"sfxuWz8KldmlBaUOAvyhQQ2e",
			"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name : "auth",
			Value : "",
			Path : "/",
			MaxAge : -1,
		})
		w.Header().Set("Location","/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	go r.run()

	// Start Web Server
	log.Println("Start Web Server. Port ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
