package bigramms

import "github.com/BaldiSlayer/rofl-lab3/internal/grammar"

type Bigramms struct {
	First   []int
	Last    []int
	Follow  []int
	Precede []int
}

func (b *Bigramms) Build(g *grammar.Grammar) *Bigramms {
	return &Bigramms{}
}
