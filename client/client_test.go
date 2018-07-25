package client

import (
	"testing"
	"fmt"
)

func TestGet(t *testing.T) {
	html, err := get(
		"http://www2.correios.com.br/sobrecorreios/empresa/acessoinformacao/servidores/ListaServidores/" +
			"lisServidores.cfm?letra=A",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(html))
}

func TestPaginas(t *testing.T) {
	letras, err := Letras()
	if err != nil {
		t.Fatal(err)
	}
	if len(letras) < 27 || len(letras) > 40 {
		t.Fatal("quantidade de letras muito fora do esperado")
	}
}

func TestFuncionarios(t *testing.T) {
	funcionarios, err := Funcionarios("Y")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range funcionarios {
		fmt.Println(f)
	}
}
