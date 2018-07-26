package server

import (
	"net/http"
	"fmt"
	"github.com/luisfernandogaido/funcionarios/modelo"
	"strings"
	"io/ioutil"
	"encoding/json"
)

func drs(w http.ResponseWriter, r *http.Request) {
	drs, err := modelo.Drs()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	printJson(w, drs)
}

func dr(w http.ResponseWriter, r *http.Request) {
	dr := strings.Replace(r.URL.Path, "/funcionarios/drs/", "", 1)
	if dr == "" {
		http.Error(w, "dr obrigatória", http.StatusBadRequest)
		return
	}
	funcionarios, err := modelo.FuncionariosDr(dr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	printJson(w, funcionarios)
}

func referencias(w http.ResponseWriter, r *http.Request) {
	referencias, err := modelo.Referencias()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	printJson(w, referencias)
}

func referencia(w http.ResponseWriter, r *http.Request) {
	referencia := strings.Replace(r.URL.Path, "/funcionarios/referencias/", "", 1)
	if referencia == "" {
		http.Error(w, "referência obrigatória", http.StatusBadRequest)
		return
	}
	funcionarios, err := modelo.FuncionariosReferencia(referencia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	printJson(w, funcionarios)
}

func matricula(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		matricula := strings.Replace(r.URL.Path, "/funcionarios/matriculas/", "", 1)
		if matricula == "" {
			http.Error(w, "matrícula obrigatória", http.StatusBadRequest)
			return
		}
		funcionario, err := modelo.FuncionarioMatricula(matricula)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if funcionario.Matricula == "" {
			http.Error(w, "matrícula não encontrada.", 404)
			return
		}
		printJson(w, funcionario)
	default:
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
}

func matriculas(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var matriculas []string
	err = json.Unmarshal(bytes, &matriculas)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mapa, err := modelo.FuncionarioMatriculas(matriculas)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	printJson(w, mapa)
}

func matriculasConc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var matriculas []string
	err = json.Unmarshal(bytes, &matriculas)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mapa, err := modelo.FuncionarioMatriculasConc(matriculas, 8)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	printJson(w, mapa)
}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "oi")
}

func matriculasSorteadas(w http.ResponseWriter, r *http.Request) {
	matriculas, err := modelo.MatriculasSorteadas()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	printJson(w, matriculas)
}
