package fuzzer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cyk"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
)

func randomFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Float64()
}

func randomItem(items []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := r.Intn(len(items))

	return items[randomIndex]
}

func randomKeyFromMap(m map[string]struct{}) string {
	if len(m) == 0 {
		return ""
	}

	a := make([]string, 0, len(m))

	for item := range m {
		a = append(a, item)
	}

	return randomItem(a)
}

type Parser interface {
	// Parse parses CFG from string s to grammar with specified startSymbol
	Parse(s string, startSymbol string) *grammar.Grammar
}

type CNFer interface {
	// ToCNF translates grammar to Chomsky Normal Form
	ToCNF(g *grammar.Grammar) *grammar.Grammar
}

type InCFG interface {
	// Check checks is word can be produced by CFG
	Check(word string) bool
}

type Fuzzer struct {
	bigramm *bigramms.Bigramms
	g       *grammar.Grammar

	cyk InCFG
}

func New(s string, p Parser, cnf CNFer, b *bigramms.Bigramms, startSymbol string) *Fuzzer {
	g := p.Parse(s, startSymbol)

	gCNF := cnf.ToCNF(g)

	return &Fuzzer{
		bigramm: b.Build(gCNF),
		cyk:     cyk.New(gCNF),
		g:       gCNF,
	}
}

func (f *Fuzzer) genString(terminals []string, breakProb float64, terminalAddingProb float64) string {
	res := randomKeyFromMap(f.bigramm.First[f.g.Start])
	lastSmb := res

	for true {
		randVal := randomFloat()

		if randomFloat() < breakProb {
			break
		}

		// add terminal
		if randVal < terminalAddingProb {
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

func (f *Fuzzer) Generate(n int, breakProb, terminalAddingProb float64) []string {
	output := make([]string, n)

	// cringe
	boolToInt := func(b bool) int {
		if b {
			return 1
		}

		return 0
	}

	terminals := f.g.ExtractTerminals()

	if len(terminals) == 0 {
		fmt.Println("there are no terminals. exiting")

		return []string{}
	}

	for i := 0; i < n; i++ {
		gennedStr := f.genString(terminals, breakProb, terminalAddingProb)

		output[i] = fmt.Sprintf("%s %d", gennedStr, boolToInt(f.cyk.Check(gennedStr)))
	}

	return output
}
