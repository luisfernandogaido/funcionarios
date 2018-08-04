package client

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/luisfernandogaido/funcionarios/modelo"
	"gopkg.in/mgo.v2"
	"log"
)

const root = "http://www2.correios.com.br/sobrecorreios/empresa/acessoinformacao/servidores"

var erPg *regexp.Regexp

type Table struct {
	XMLName xml.Name `xml:"table"`
	Trs     []Tr     `xml:"tr"`
}

type Tr struct {
	Tds []string `xml:"td"`
}

type Funcionario struct {
	Nome          string
	Matricula     string
	Cpf           string
	Admissao      time.Time
	Cargo         string
	Funcao        string
	Especialidade string
	Dr            string
	Lotacao       string
	Jornada       int
	Referencia    string
	Afastamento   string
}

func init() {
	erPg = regexp.MustCompile(
		`href="sobrecorreios/empresa/acessoinformacao/servidores/ListaServidores/lisServidores.cfm\?letra=([^"]+)"`,
	)
}

func get(u string) (string, error) {
	res, err := http.Get(u)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Letras() ([]string, error) {
	html, err := get(root + "/default.cfm")
	if err != nil {
		return nil, err
	}
	matches := erPg.FindAllStringSubmatch(html, -1)
	paginas := make([]string, 0, 29)
	for _, m := range matches {
		paginas = append(paginas, m[1])
	}
	return paginas, nil
}

func Funcionarios(letra string) ([]Funcionario, error) {
	html, err := get(root + "/ListaServidores/lisServidores.cfm?letra=" + letra)
	if err != nil {
		return nil, err
	}
	newHtml := make([]rune, 0, len(html))
	for _, r := range html {
		if r == utf8.RuneError {
			continue
		}
		newHtml = append(newHtml, r)
	}
	html = string(newHtml)
	p, q := strings.Index(html, "<table>"), strings.Index(html, "</table>")
	html = html[p : q+8]
	var table Table
	if err := xml.Unmarshal([]byte(html), &table); err != nil {
		return nil, err
	}
	funcionarios := make([]Funcionario, 0)
	table.Trs = table.Trs[1:]
	for _, tr := range table.Trs {
		f := Funcionario{}
		f.Nome = strings.ToUpper(tr.Tds[0])
		f.Matricula = tr.Tds[1]
		f.Cpf = tr.Tds[2]
		admissao, err := time.Parse("02/01/2006", tr.Tds[3])
		if err == nil {
			f.Admissao = admissao
		}
		cargoFuncao := strings.Split(strings.ToUpper(tr.Tds[4]), " / ")
		f.Cargo, f.Funcao = cargoFuncao[0], cargoFuncao[1]
		f.Especialidade = strings.ToUpper(tr.Tds[5])
		drLotacao := strings.Split(tr.Tds[6], " / ")
		f.Dr, f.Lotacao = drLotacao[0], drLotacao[1]
		jornada, err := strconv.Atoi(strings.Replace(tr.Tds[7], "h", "", -1))
		if err == nil {
			f.Jornada = jornada
		}
		f.Jornada = jornada
		f.Referencia = tr.Tds[8]
		f.Afastamento = tr.Tds[9]
		funcionarios = append(funcionarios, f)
	}
	return funcionarios, nil
}

func Importa() error {
	letras, err := Letras()
	if err != nil {
		return err
	}
	conc := 29
	sem := make(chan struct{}, conc)
	chErr := make(chan error)
	for _, letra := range letras {
		sem <- struct{}{}
		go func(l string) {
			defer func() { <-sem }()
			funcionarios, err := Funcionarios(l)
			if err != nil {
				chErr <- err
				return
			}
			for _, f := range funcionarios {
				fun := modelo.Funcionario{}
				fun.Matricula = f.Matricula
				fun.Nome = f.Nome
				fun.Cpf = f.Cpf
				fun.Admissao = f.Admissao
				fun.Cargo = f.Cargo
				fun.Funcao = f.Funcao
				fun.Especialidade = f.Especialidade
				fun.Dr = f.Dr
				fun.Lotacao = f.Lotacao
				fun.Referencia = f.Referencia
				fun.Afastamento = f.Afastamento
				fun.Indice = f.Matricula + " " + f.Nome + " " + f.Cargo + " " + f.Funcao + " " + f.Especialidade + " " +
					f.Dr + " " + f.Lotacao + " " + f.Referencia
				if err := fun.Salva(); err != nil {
					chErr <- err
				}
			}
		}(letra)
	}
	terminados := 0
	for {
		select {
		case sem <- struct{}{}:
			terminados++
			if terminados == conc {
				return nil
			}
		case err := <-chErr:
			fmt.Println(err)
		}
	}
}

func ImportaMgo() error {

	//sess, err := mgo.Dial("127.0.0.1:27017")
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"104.131.64.134:27017"},
		Timeout:  60 * time.Second,
		Database: "admin",
		Username: "root",
		Password: "1000sonhosreais",
	})
	if err != nil {
		log.Fatal(err)
	}
	collFun := sess.DB("funcionarios").C("funcionarios")
	letras, err := Letras()
	if err != nil {
		return err
	}
	conc := 29
	sem := make(chan struct{}, conc)
	chErr := make(chan error)
	for _, letra := range letras {
		sem <- struct{}{}
		go func(l string) {
			defer func() { <-sem }()
			funcionarios, err := Funcionarios(l)
			if err != nil {
				chErr <- err
				return
			}
			inter := make([]interface{}, 0, len(funcionarios))
			for _, f := range funcionarios {
				inter = append(inter, f)
				continue
				fun := modelo.Funcionario{}
				fun.Matricula = f.Matricula
				fun.Nome = f.Nome
				fun.Cpf = f.Cpf
				fun.Admissao = f.Admissao
				fun.Cargo = f.Cargo
				fun.Funcao = f.Funcao
				fun.Especialidade = f.Especialidade
				fun.Dr = f.Dr
				fun.Lotacao = f.Lotacao
				fun.Referencia = f.Referencia
				fun.Afastamento = f.Afastamento
				fun.Indice = f.Matricula + " " + f.Nome + " " + f.Cargo + " " + f.Funcao + " " + f.Especialidade + " " +
					f.Dr + " " + f.Lotacao + " " + f.Referencia
				if err := collFun.Insert(fun); err != nil {
					chErr <- err
				}
			}
			if err := collFun.Insert(inter...); err != nil {
				chErr <- err
			}
		}(letra)
	}
	terminados := 0
	for {
		select {
		case sem <- struct{}{}:
			terminados++
			if terminados == conc {
				return nil
			}
		case err := <-chErr:
			fmt.Println(err)
		}
	}
}
