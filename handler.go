package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}

func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "hello world")
}

func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Rate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
