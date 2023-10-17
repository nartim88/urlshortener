package main

import (
	"log"
	"net/http"
)

type custURL struct {
	full    string
	shorten string
}

type store struct {
	custURL
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
