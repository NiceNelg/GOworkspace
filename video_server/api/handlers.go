package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CreaterUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(w, "CreaterUser")
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username := p.ByName("username")
	io.WriteString(w, username)
}
