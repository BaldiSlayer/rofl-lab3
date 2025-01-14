package app

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cyk"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
	"math/rand"
	"time"
)

type Parser interface {
	Parse(s string) *grammar.Grammar
}

type CNFer interface {
	ToCNF(g *grammar.Grammar) *grammar.Grammar
}

type Fuzzer struct {
	bigramm *bigramms.Bigramms
	cyk     *cyk.CYK
}

func New(p parser.Parser, cnf CNFer, b bigramms.Bigramms) *Fuzzer {
	s := ""

	gram := p.Parse(s)

	gCNF := cnf.ToCNF(gram)

	bm := b.Build(gCNF)
	c := cyk.New(gCNF)

	return &Fuzzer{
		bigramm: bm,
		cyk:     c,
	}
}

func randSliceItem(a []int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(len(a))
}

// lowercaseLetters := "abcdefghijklmnopqrstuvwxyz"

func (f *Fuzzer) Generate(n int) {
	sTerminals := []int{1, 2}

	for i := 0; i < n; i++ {
		_ = randSliceItem(sTerminals)

	}
}
