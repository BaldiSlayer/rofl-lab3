package fuzzer

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cyk"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
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
	g       *grammar.Grammar
}

func New(p Parser, cnf CNFer, b *bigramms.Bigramms) *Fuzzer {
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

// lowercaseLetters := "abcdefghijklmnopqrstuvwxyz"

func randomKey(m map[string]map[string]struct{}) string {
	if len(m) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(keys))

	return keys[randomIndex]
}

func randomFloat() float64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Float64()
}

func randomItem(items []string) string {
	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(items))

	return items[randomIndex]
}

func randomKeyFromMap(m map[string]struct{}) string {
	a := make([]string, len(m))

	for item := range m {
		a = append(a, item)
	}

	return randomItem(a)
}

func (f *Fuzzer) genString(terminals []string, someValue float64) string {
	res := randomKey(f.bigramm.First)
	lastSmb := res

	for true {
		randVal := randomFloat()

		if randomFloat() < 0.1 {
			break
		}

		// add terminal
		if randVal < someValue {
			lastSmb = randomItem(terminals)
			res += lastSmb

			continue
		}

		_, ok := f.bigramm.Matrix[lastSmb]
		cond := ok && len(f.bigramm.Matrix[lastSmb]) != 0

		if !cond {
			break
		}

		lastSmb = randomKeyFromMap(f.bigramm.Matrix[string(res[len(res)-1])])
		res += lastSmb
	}

	return res
}

func (f *Fuzzer) Generate(n int, startSmb string) []string {
	output := make([]string, 0, n)

	terminals := f.g.ExtractTerminals()

	for i := 0; i < n; i++ {
		gennedStr := f.genString(terminals, 0.1)

		belongs := f.cyk.Check(gennedStr)

		output = append(output, fmt.Sprintf("%s %t", gennedStr, belongs))
	}

	return output
}
