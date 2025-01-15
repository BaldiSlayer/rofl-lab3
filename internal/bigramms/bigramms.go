package bigramms

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
)

type Bigramms struct {
	Matrix map[string]map[string]struct{}
}

func isNotTerminal(symbols models.SymbolsBtw) bool {
	return !(symbols.S[0] >= 'a' && symbols.S[0] <= 'z')
}

// todo is terminal must be normal
func isTerminal(symbols models.SymbolsBtw) bool {
	return !isNotTerminal(symbols)
}

func union(dest map[string]struct{}, src map[string]struct{}) map[string]struct{} {
	for key := range src {
		dest[key] = struct{}{}
	}

	return dest
}

func makeFirstAndLastRec(
	g *grammar.Grammar,
	visited map[string]struct{},
	nt string,
	first map[string]map[string]struct{},
	last map[string]map[string]struct{},
) {
	visited[nt] = struct{}{}

	// todo not cool function name
	step := func(targetSet map[string]map[string]struct{}, smb models.SymbolsBtw) {
		if isNotTerminal(smb) {
			if _, ok := visited[smb.S]; !ok {
				makeFirstAndLastRec(g, visited, smb.S, targetSet, last)
			}

			targetSet[nt] = union(targetSet[nt], targetSet[smb.S])
		} else {
			targetSet[nt][smb.S] = struct{}{}
		}
	}

	for _, rightRule := range g.Grammar[nt].Rights {
		// update first
		step(first, rightRule.Body[0])

		// update last
		step(first, rightRule.Body[len(rightRule.Body)-1])
	}
}

func makeFirstAndLast(g *grammar.Grammar) (map[string]map[string]struct{}, map[string]map[string]struct{}) {
	visited := make(map[string]struct{})
	first := make(map[string]map[string]struct{})
	last := make(map[string]map[string]struct{})

	for nt := range g.Grammar {
		if _, ok := visited[nt]; !ok {
			makeFirstAndLastRec(g, visited, nt, first, last)
		}
	}

	return first, last
}

func checkFollow(
	g *grammar.Grammar,
	follow map[string]map[string]struct{},
	first map[string]map[string]struct{},
) (map[string]map[string]struct{}, bool) {
	changed := false

	isChanged := func(rightRule models.ProductionBody) bool {
		for terminal := range first[rightRule.Body[1].S] {
			if _, ok := follow[rightRule.Body[0].S][terminal]; !ok {
				return true
			}
		}

		return false
	}

	for _, rightRules := range g.Grammar {
		for _, rightRule := range rightRules.Rights {
			if len(rightRule.Body) > 1 {
				changed = isChanged(rightRule)

				follow[rightRule.Body[0].S] = union(
					follow[rightRule.Body[0].S],
					first[rightRule.Body[1].S],
				)
			}
		}
	}

	return follow, changed
}

func makeFollow(g *grammar.Grammar, first map[string]map[string]struct{}) map[string]map[string]struct{} {
	follow := make(map[string]map[string]struct{})

	var changed bool

	for {
		follow, changed = checkFollow(g, follow, first)

		if !changed {
			return follow
		}
	}
}

func checkPrecede(
	g *grammar.Grammar,
	precede map[string]map[string]struct{},
	last map[string]map[string]struct{},
) (map[string]map[string]struct{}, bool) {
	changed := false

	isChanged := func(rightRule models.ProductionBody) bool {
		for terminal := range last[rightRule.Body[0].S] {
			if _, ok := precede[rightRule.Body[1].S][terminal]; !ok {
				return true
			}
		}

		return false
	}

	for _, rightRules := range g.Grammar {
		for _, rightRule := range rightRules.Rights {
			if len(rightRule.Body) > 1 {
				changed = isChanged(rightRule)

				precede[rightRule.Body[1].S] = union(
					precede[rightRule.Body[1].S],
					last[rightRule.Body[0].S],
				)
			}
		}
	}

	return precede, changed
}

func makePrecede(g *grammar.Grammar, first map[string]map[string]struct{}) map[string]map[string]struct{} {
	precede := make(map[string]map[string]struct{})

	var changed bool

	for {
		precede, changed = checkPrecede(g, precede, first)

		if !changed {
			return precede
		}
	}
}

func extractTerminalsFromRule(rule models.ProductionBody) []string {
	terminals := make([]string, 0)

	for _, smb := range rule.Body {
		if isTerminal(smb) {
			terminals = append(terminals, smb.S)
		}
	}

	return terminals
}

func extractTerminalsFromGrammar(g *grammar.Grammar) []string {
	terminals := make([]string, 0, len(g.Grammar))

	for _, rules := range g.Grammar {
		for _, rule := range rules.Rights {
			terminals = append(terminals, extractTerminalsFromRule(rule)...)
		}
	}

	return terminals
}

// very bad function
func needToAdd(
	y1, y2 string,
	first, last, follow, precede map[string]map[string]struct{},
) bool {
	for a1 := range last {
		_, ok1 := last[a1][y1]
		_, ok2 := follow[a1][y2]

		if ok1 && ok2 {
			return true
		}
	}

	for a1 := range precede {
		_, ok1 := precede[a1][y1]
		_, ok2 := first[a1][y2]

		if ok1 && ok2 {
			return true
		}
	}

	for a1 := range last {
		_, ok1 := last[a1][y1]
		_, ok2 := first[a1][y1]
		_, ok3 := follow[a1][y2]

		if ok1 && ok2 && ok3 {
			return true
		}
	}

	return false
}

func pairChecking(g *grammar.Grammar) func(y1, y2 string) bool {
	exists := make(map[string]struct{})

	for _, rules := range g.Grammar {
		for _, rule := range rules.Rights {
			if len(rule.Body) == 2 {
				exists[fmt.Sprintf("%s %s", rule.Body[0].S, rule.Body[1].S)] = struct{}{}
			}
		}
	}

	return func(y1, y2 string) bool {
		_, ok := exists[y1+" "+y2]

		return ok
	}
}

func makeBigramMatrix(
	g *grammar.Grammar,
	first, last, follow, precede map[string]map[string]struct{},
) map[string]map[string]struct{} {
	bigramms := make(map[string]map[string]struct{})

	pairChecker := pairChecking(g)
	terminals := extractTerminalsFromGrammar(g)

	// am i need to check that y1 != y2?
	for _, y1 := range terminals {
		for _, y2 := range terminals {
			if pairChecker(y1, y2) || needToAdd(y1, y2, first, last, follow, precede) {
				if _, ok := bigramms[y1]; !ok {
					bigramms[y1] = make(map[string]struct{})
				}

				bigramms[y1][y2] = struct{}{}
			}
		}
	}

	return bigramms
}

func (b *Bigramms) Build(g *grammar.Grammar) *Bigramms {
	first, last := makeFirstAndLast(g)
	follow := makeFollow(g, first)
	precede := makePrecede(g, last)

	matrix := makeBigramMatrix(g, first, last, follow, precede)

	return &Bigramms{
		Matrix: matrix,
	}
}
