package helper

import "testing"

func TestTemplatesParse(t *testing.T) {
	r := InitRenderer("../views")
	if r == nil || len(r.templates) == 0 {
		t.Fatal("aucune vue chargee")
	}
}
