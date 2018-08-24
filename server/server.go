package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Start(addr string) error {
	http.HandleFunc("/funcionarios/drs", drs)
	http.HandleFunc("/funcionariosmongo/drs", drsMongo)
	http.HandleFunc("/funcionarios/drs/", dr)
	http.HandleFunc("/funcionariosmongo/drs/", drMongo)
	http.HandleFunc("/funcionarios/referencias", referencias)
	http.HandleFunc("/funcionariosmongo/referencias", referenciasMongo)
	http.HandleFunc("/funcionarios/referencias/", referencia)
	http.HandleFunc("/funcionarios/matriculas", matriculas)
	http.HandleFunc("/funcionarios/matriculasconc", matriculasConc)
	http.HandleFunc("/funcionarios/matriculas/", matricula)
	http.HandleFunc("/funcionarios/search", search)
	http.HandleFunc("/matriculas/sorteadas", matriculasSorteadas)
	fmt.Printf("Ouvindo em %v...\n", addr)
	return http.ListenAndServe(addr, nil)
}

func printJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Add("Content-type", "application/json; charset=utf8")
	bytes, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(bytes))
	return err
}

func cors(w http.ResponseWriter) {
	allowedHeaders := "Content-type, Cache-Control"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
}
