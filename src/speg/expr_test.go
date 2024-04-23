package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestExprParser(t *testing.T) {

	varname := Token(Letters()).Tagged("var")
	number := Token(Digits()).Tagged("num")
	lparen := Token(Exactly("("))
	rparen := Token(Exactly(")"))
	addOp := Token(Or(Exactly("+"), Exactly("-")))
	multOp := Token(Or(Exactly("*"), Exactly("/")))
	var expr Parser
	
	factor := Or(
		varname,
		number,
		Seq(lparen.Omit(),
			Indirect(&expr),
			rparen.Omit(),
		).Tagged("expr"),
	)

	term := Left(
		factor,
		Seq(multOp, factor),
	).Tagged("prod")

	expr = Left(
		term,
		Seq(addOp, term),
	).Tagged("sum")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"base", "x", `(var "x")`},
		{"add", "x+y", `(sum (var "x") ("+") (var "y"))`},
		{"add", "x + y", `(sum (var "x") ("+") (var "y"))`},
		{"mult", "x*y", `(prod (var "x") ("*") (var "y"))`},
		{"left associative", "x+y+z", `(sum (sum (var "x") ("+") (var "y")) ("+") (var "z"))`},
		{"left assoc neg", "x+y-z", `(sum (sum (var "x") ("+") (var "y")) ("-") (var "z"))`},
		{"precedence", "x+y*z", `(sum (var "x") ("+") (prod (var "y") ("*") (var "z")))`},
		{"ditto", "x+y*3", `(sum (var "x") ("+") (prod (var "y") ("*") (num "3")))`},
		{"parens", "x+y*(3+2)", `(sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2")))))`},
		{"precedence with parens", "x+y*(3+2)", `(sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2")))))`},
		{"complicated", "x+y*(3+2)+(7*3)", `(sum (sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2"))))) ("+") (expr (prod (num "7") ("*") (num "3"))))`},
		{"spaces", "x + y * (3+2) + (7*3)", `(sum (sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2"))))) ("+") (expr (prod (num "7") ("*") (num "3"))))`},
		{"spaces2", " x + y * ( 3 + 2 ) + (7*3)", `(sum (sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2"))))) ("+") (expr (prod (num "7") ("*") (num "3"))))`},
		{"spaces2", "	x 	+ y * ( 3 + 2 ) + (7*3)", `(sum (sum (var "x") ("+") (prod (var "y") ("*") (expr (sum (num "3") ("+") (num "2"))))) ("+") (expr (prod (num "7") ("*") (num "3"))))`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, expr.Parse([]rune(tc.input), 0, NewContext()).String())
		})
	}
}
