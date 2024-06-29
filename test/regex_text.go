package test

import (
	"fmt"
	"regex_in_go/utils"
	"testing"
)

func TestNfa(t *testing.T) {

	var data = []struct {
		email    string
		validity bool
	}{
		// Salidas Salidas no validas cuando se compara con la expresion deseada
		{email: "valid_email@example.com", validity: true},
		{email: "john.doe@email.com", validity: true},
		{email: "user_name@email.org", validity: true},
		{email: "support@email.io", validity: true},
		{email: "contact@123.com", validity: true},
		{email: "sales@email.biz", validity: true},

		// Salidas no validas cuando se compara con la expresion deseada
		{email: "alice.smith123@email.co.uk", validity: false},
		{email: "invalid.email@", validity: false},
		{email: ".invalid@email.com", validity: false},
		{email: "email@invalid..com", validity: false},
		{email: "user@-invalid.com", validity: false},
		{email: "user@invalid-.com", validity: false},
	}

	ctx := utils.Parse(`[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}`)
	nfa := utils.ToNfa(ctx)

	for _, instance := range data {
		t.Run(fmt.Sprintf("Test: '%s'", instance.email), func(t *testing.T) {
			result := nfa.Check(instance.email, -1)
			if result != instance.validity {
				t.Logf("Expected: %t, got: %t\n", instance.validity, result)
				t.Fail()
			}
		})
	}
}
