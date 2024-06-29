package main

import (
	"fmt"
	"regex_in_go/utils"
)

func main() {

	var data = []struct {
		email string
	}{
		{email: "valid_email@example.com"},
		{email: "user@sub.domain.1a"}}

	regex := utils.Parse(`[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}`)
	nfa := utils.ToNfa(regex)

	for _, i := range data {
		res := nfa.Check(i.email, -1)
		if !res {
			fmt.Println("No se cumple la condicion")
		} else {
			fmt.Println("Se cumple la condicion")
		}
	}
}
