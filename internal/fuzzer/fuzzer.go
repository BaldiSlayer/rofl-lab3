package app

import "github.com/BaldiSlayer/rofl-lab3/internal/grammar"

type Parser interface {
	Parse(s string) *grammar.Grammar
}

type CNFer interface {
	ToCNF(g *grammar.Grammar) *grammar.Grammar
}

type Fuzzer struct {
	parser Parser
	cnf    CNFer
}

func New() *Fuzzer {
	return &Fuzzer{}
}

func (f *Fuzzer) Run() {
	s := ""

	gram := f.parser.Parse(s)

	_ = f.cnf.ToCNF(gram)
}
