package main

import (
	"log"
	"fmt"
	"time"

	"github.com/luisfernandogaido/funcionarios/client"
)

func main() {
	t0 := time.Now()
	letras, err := client.Letras()
	if err != nil {
		log.Fatal(err)
	}
	concorrencia := 32
	chInt := make(chan int)
	sem := make(chan struct{}, concorrencia)
	total := 0
	go func() {
		for n := range chInt {
			total += n
		}
	}()
	for _, letra := range letras {
		sem <- struct{}{}
		go func(letra string) {
			defer func() { <-sem }()
			funcionarios, err := client.Funcionarios(letra)
			if err != nil {
				log.Fatal(err)
			}
			chInt <- len(funcionarios)
			fmt.Println(letra, len(funcionarios))
		}(letra)
	}
	for i := 0; i < concorrencia; i++ {
		sem <- struct{}{}
	}
	time.Sleep(time.Millisecond)
	close(chInt)
	fmt.Println(total)
	fmt.Println(time.Since(t0))
}
