package main

import (
	"log"
	"net/http"
)

var SavedData map[string]string

type CustURL struct {
	full    string
	shorten string
}

type Repository interface {
	Save() bool
	IsExist() bool
}

func (url *CustURL) Save() bool {
	if url.IsExist() {
		return false
	}
	SavedData[url.shorten] = url.full
	return true
}

func (url *CustURL) IsExist() bool {
	if _, ok := SavedData[url.shorten]; !ok {
		return false
	}
	return true
}

var allowedMethods = map[string]string{
	"GET":  http.MethodGet,
	"POST": http.MethodPost,
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	meth, ok := allowedMethods[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed."))
	}
	if meth == http.MethodPost {
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
