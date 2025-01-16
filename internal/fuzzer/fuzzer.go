package fuzzer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cyk"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
)

type Parser interface {
	Parse(s string, startSymbol string) *grammar.Grammar
}

type CNFer interface {
	ToCNF(g *grammar.Grammar) *grammar.Grammar
}

type InCFG interface {
	Check(word string) bool
}

type Fuzzer struct {
	bigramm *bigramms.Bigramms
	g       *grammar.Grammar

	cyk InCFG
}

func New(s string, p Parser, cnf CNFer, b *bigramms.Bigramms, startSymbol string) *Fuzzer {
	gram := p.Parse(s, startSymbol)

	gCNF := cnf.ToCNF(gram)

	bm := b.Build(gCNF)
	c := cyk.New(gCNF)

	return &Fuzzer{
		bigramm: bm,
		cyk:     c,
		g:       gCNF,
	}
}

func randomFloat() float64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Float64()
}

func randomItem(items []string) string {
	rand.Seed(time.Now().UnixNano())

	if len(items) == 0 {
		return ""
	}

	randomIndex := rand.Intn(len(items))

	return items[randomIndex]
}

func randomKeyFromMap(m map[string]struct{}) string {
	a := make([]string, 0, len(m))

	for item := range m {
		a = append(a, item)
	}

	return randomItem(a)
}

func (f *Fuzzer) genString(terminals []string, someValue float64, startSmb string) string {
	res := randomKeyFromMap(f.bigramm.First[startSmb])
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

		// add from bigrams
		lastSmb = randomKeyFromMap(f.bigramm.Matrix[string(res[len(res)-1])])
		res += lastSmb
	}

	return res
}

// cringe
func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

func (f *Fuzzer) Generate(n int, someValue float64, startSmb string) []string {
	//fmt.Println(f.g.Print())

	output := make([]string, 0, n)

	terminals := f.g.ExtractTerminals()

	for i := 0; i < n; i++ {
		gennedStr := f.genString(terminals, someValue, startSmb)

		output = append(
			output,
			fmt.Sprintf("%s %d", gennedStr, boolToInt(f.cyk.Check(gennedStr))),
		)
	}

	return output
}
