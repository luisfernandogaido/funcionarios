package client

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"encoding/xml"
	"time"
	"unicode/utf8"
	"strings"
	"strconv"
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
