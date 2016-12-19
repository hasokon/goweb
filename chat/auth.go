package main

import (
	"strings"
	"net/http"
	"fmt"
	"log"
)
//CLientID
//176223784373-bss719fm64msnfuqau908ql7l1es0dcs.apps.googleusercontent.com
//Client Secret
//sfxuWz8KldmlBaUOAvyhQQ2e

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//Uncertified
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		//Some kind of error
		panic(err.Error())
	} else {
		//Success
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
		case "login" :
			log.Println("TODO: Login Processing", provider)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Action '%s' is not suported", action)
	}
}
