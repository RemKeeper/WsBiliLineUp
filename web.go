package main

import (
	"net/http"
)

func WebServerTwo() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", indexHandler)
	serverMux.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("./"))))
	http.ListenAndServe(":100", serverMux)
}
