package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func Test(t *testing.T) {
	var expr Parser
	var iexpr = Indirect()

	iexpr.Set(&expr)
	varname := Letters().Tagged("var")
	number := Digits().Tagged("num")
	factor := Or(
		varname,
		number,
		Seq(Exactly("(").Omit(),
			iexpr,
			Exactly(")").Omit(),
		).Tagged("expr"),
	)
	term := Left(
		factor,
		Seq(
			Or(Exactly("*"), Exactly("/")),
			factor,
		),
	).Tagged("prod")
	expr = Left(
		term,
		Seq(
			Or(
				Exactly("+"),
				Exactly("-"),
			),
			term,
		)).Tagged("sum")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"base", "x", `(var "x")`},
		{"base", "x+y", `(sum (var "x") "+" (var "y"))`},
		{"base", "x*y", `(prod (var "x") "*" (var "y"))`},
		{"base", "x+y+z", `(sum (sum (var "x") "+" (var "y")) "+" (var "z"))`},
		{"base", "x+y-z", `(sum (sum (var "x") "+" (var "y")) "-" (var "z"))`},
		{"base", "x+y*z", `(sum (var "x") "+" (prod (var "y") "*" (var "z")))`},
		{"base", "x+y*3", `(sum (var "x") "+" (prod (var "y") "*" (num "3")))`},
		{"base", "x+y*(3+2)", `(sum (var "x") "+" (prod (var "y") "*" (expr (sum (num "3") "+" (num "2")))))`},
		{"base", "x+y*(3+2)", `(sum (var "x") "+" (prod (var "y") "*" (expr (sum (num "3") "+" (num "2")))))`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, expr.Parse([]rune(tc.input), 0, NewContext()).String())
		})
	}
}
