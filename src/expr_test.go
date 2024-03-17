package sparse

import (
	"github.com/shoenig/test"
	"testing"
)

func TestExpr(t *testing.T) {
	var Term, Expr Parser

	Term = FirstOf(
		Digits,
		Seq(Deref(&Term), Exactly("*"), Deref(&Term)),
		Seq(Exactly("("), Deref(&Expr), Exactly(")")),
	)

	Expr = FirstOf(
		Seq(Term, Exactly("+"), Term),
		Term,
	)

	test.Eq(t, "1+2", Expr([]rune("1+2*3")).String())
}

func TestAlt(t *testing.T) {
	p := Longest(Alt(Letters, Seq(Letters, Digits), Letters))
	test.Eq(t, "abc123", p([]rune("abc123")).String())

}
