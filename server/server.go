package server

import (
	"fmt"
	"net/http"
)

func Start(addr string) error {
	http.HandleFunc("/funcionarios/search", search)
	return http.ListenAndServe(addr, nil)
}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "oi")
}
