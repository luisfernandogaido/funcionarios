package modelo

import (
	"time"
	"strings"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"encoding/json"
)

type Funcionario struct {
	Matricula     string    `json:"matricula"`
	Nome          string    `json:"nome"`
	Cpf           string    `json:"cpf"`
	Admissao      time.Time `json:"admissao"`
	Cargo         string    `json:"cargo"`
	Funcao        string    `json:"funcao"`
	Especialidade string    `json:"especialidade"`
	Dr            string    `json:"dr"`
	Lotacao       string    `json:"lotacao"`
	Jornada       int       `json:"jornada"`
	Referencia    string    `json:"referencia"`
	Afastamento   string    `json:"afastamento"`
	Indice        string    `json:"indice"`
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

type drRes struct {
	Dr          string `json:"dr"`
	Ocorrencias int    `json:"ocorrencias"`
}

func Drs() ([]drRes, error) {
	rows, err := db.Query("call drs_seleciona()")
	if err != nil {
		return nil, err
	}
	drs := make([]drRes, 0, 29)
	for rows.Next() {
		dr := drRes{}
		err = rows.Scan(&dr.Dr, &dr.Ocorrencias)
		if err != nil {
			return nil, err
		}
		dr.Dr = strings.Replace(dr.Dr, "SE/", "", 1)
		drs = append(drs, dr)
	}
	return drs, nil

}

func FuncionariosDr(dr string) ([]Funcionario, error) {
	dr = "SE/" + dr
	rows, err := db.Query("call funcionarios_dr_seleciona(?)", dr)
	if err != nil {
		return nil, err
	}
	ff := make([]Funcionario, 0)
	for rows.Next() {
		f := Funcionario{}
		err = rows.Scan(
			&f.Matricula,
			&f.Nome,
			&f.Cpf,
			&f.Admissao,
			&f.Cargo,
			&f.Funcao,
			&f.Especialidade,
			&f.Dr,
			&f.Lotacao,
			&f.Jornada,
			&f.Referencia,
			&f.Afastamento,
			&f.Indice,
		)
		if err != nil {
			return nil, err
		}
		if f.Referencia == "" {
			f.Referencia = "NULL"
		}
		ff = append(ff, f)
	}
	return ff, nil
}

type referenciaRes struct {
	Referencia  string `json:"referencia"`
	Ocorrencias int    `json:"ocorrencias"`
}

func Referencias() ([]referenciaRes, error) {
	rows, err := db.Query("call referencias_seleciona()")
	if err != nil {
		return nil, err
	}
	referencias := make([]referenciaRes, 0, 150)
	for rows.Next() {
		ref := referenciaRes{}
		err = rows.Scan(&ref.Referencia, &ref.Ocorrencias)
		if err != nil {
			return nil, err
		}
		if ref.Referencia == "" {
			ref.Referencia = "NULL"
		}
		referencias = append(referencias, ref)
	}
	return referencias, nil
}

func FuncionariosReferencia(referencia string) ([]Funcionario, error) {
	if referencia == "NULL" {
		referencia = ""
	}
	rows, err := db.Query("call funcionarios_referencia_seleciona(?)", referencia)
	if err != nil {
		return nil, err
	}
	ff := make([]Funcionario, 0)
	for rows.Next() {
		f := Funcionario{}
		err = rows.Scan(
			&f.Matricula,
			&f.Nome,
			&f.Cpf,
			&f.Admissao,
			&f.Cargo,
			&f.Funcao,
			&f.Especialidade,
			&f.Dr,
			&f.Lotacao,
			&f.Jornada,
			&f.Referencia,
			&f.Afastamento,
			&f.Indice,
		)
		if err != nil {
			return nil, err
		}
		if f.Referencia == "" {
			f.Referencia = "NULL"
		}
		ff = append(ff, f)
	}
	return ff, nil
}

func FuncionarioMatricula(matricula string) (Funcionario, error) {
	var fun Funcionario
	muRd.Lock()
	existe, err := redis.Bool(rd.Do("EXISTS", matricula))
	muRd.Unlock()
	if err != nil {
		return Funcionario{}, err
	}
	if existe {
		bytes, err := redis.Bytes(rd.Do("GET", matricula))
		if err != nil {
			return Funcionario{}, err
		}
		if err := json.Unmarshal(bytes, &fun); err != nil {
			return Funcionario{}, err
		}
		return fun, nil
	}
	err = db.QueryRow("CALL funcionario_seleciona(?)", matricula).Scan(
		&fun.Matricula,
		&fun.Nome,
		&fun.Cpf,
		&fun.Admissao,
		&fun.Cargo,
		&fun.Funcao,
		&fun.Especialidade,
		&fun.Dr,
		&fun.Lotacao,
		&fun.Jornada,
		&fun.Referencia,
		&fun.Afastamento,
		&fun.Indice,
	)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return Funcionario{}, err
	}

	bytes, err := json.Marshal(fun)
	if err != nil {
		return Funcionario{}, err
	}
	muRd.Lock()
	_, err = rd.Do("SET", matricula, bytes, "EX", 86400)
	muRd.Unlock()
	return fun, err
}

func FuncionarioMatriculas(matriculas []string) (map[string]Funcionario, error) {
	mapa := make(map[string]Funcionario)
	for _, m := range matriculas {
		fun, err := FuncionarioMatricula(m)
		if err != nil {
			return nil, err
		}
		mapa[m] = fun
	}
	return mapa, nil
}

func FuncionarioMatriculasConc(matriculas []string, concorrencia int) (map[string]Funcionario, error) {
	mapa := make(map[string]Funcionario)
	for _, m := range matriculas {
		fun, err := FuncionarioMatricula(m)
		if err != nil {
			return nil, err
		}
		mapa[m] = fun
	}
	return mapa, nil
}

func MatriculasSorteadas() ([]string, error) {
	rows, err := db.Query("CALL matriculas_sorteadas_seleciona()")
	if err != nil {
		return nil, err
	}
	matriculas := make([]string, 0, 50)
	for rows.Next() {
		var matricula string
		err = rows.Scan(&matricula)
		if err != nil {
			return nil, err
		}
		matriculas = append(matriculas, matricula)
	}
	return matriculas, nil
}
