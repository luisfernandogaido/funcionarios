package modelo

import "time"

type Funcionario struct {
	Matricula     string
	Nome          string
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
	Indice        string
}

func (f Funcionario) Salva() error {
	_, err := db.Exec(
		"call funcionario_insere(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		f.Matricula,
		f.Nome,
		f.Cpf,
		f.Admissao,
		f.Cargo,
		f.Funcao,
		f.Especialidade,
		f.Dr,
		f.Lotacao,
		f.Jornada,
		f.Referencia,
		f.Afastamento,
		f.Indice,
	)
	return err
}
