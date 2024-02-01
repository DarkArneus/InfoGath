package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
)

func detectInput(domains []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, domain := range domains {
		// Crea un nuevo colector
		c := colly.NewCollector()
		// Define lo que se debe hacer cuando se encuentran elementos input en la página
		c.OnHTML("input", func(e *colly.HTMLElement) {
			inputType := e.Attr("type")
			inputName := e.Attr("name")
			results <- fmt.Sprintf("Dominio: %s, Tipo de entrada: %s, Nombre: %s\n", domain, inputType, inputName)
		})
		// Define lo que se debe hacer cuando se encuentran elementos form en la página
		c.OnHTML("form", func(e *colly.HTMLElement) {
			formAction := e.Attr("action")
			results <- fmt.Sprintf("Dominio: %s, Formulario encontrado con acción: %s\n", domain, formAction)
		})
		// Visita el dominio actual
		err := c.Visit(domain)
		if err != nil {
			fmt.Printf("Error al visitar el dominio %s: %v\n", domain, err)
		}
	}
}
