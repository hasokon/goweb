package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
)

func main() {
	addr := flag.String("addr", ":8080", "Apprication's Address")
	flag.Parse()
	//Gomniauth Set-up
	gomniauth.SetSecurityKey("dGhKIuhkmhtqxNNbUuiJNVtgun[;?<,miHNbgojj")
	gomniauth.WithProviders(
		facebook.New("","","http://localhost:8080/auth/callback/facebook"),
		github.New(
			"aec2bc42117b12eb1328",
			"56282cd7e64284f019554221fad29d12aca52714",
			"http://localhost:8080/auth/callback/github"),
		google.New(
		"176223784373-bss719fm64msnfuqau908ql7l1es0dcs.apps.googleusercontent.com",
		"sfxuWz8KldmlBaUOAvyhQQ2e",
		"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	// Start Web Server
	log.Println("Start Web Server. Port ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
