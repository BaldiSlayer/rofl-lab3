package bigramms

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
)

type Bigramms struct {
	Matrix map[string]map[string]struct{}
	First  map[string]map[string]struct{}
}

func union(dest map[string]struct{}, src map[string]struct{}) map[string]struct{} {
	if dest == nil {
		dest = make(map[string]struct{})
	}

	for key := range src {
		dest[key] = struct{}{}
	}

	return dest
}

func difference(setA, setB map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})

	// Добавляем все элементы из setA в результат
	for key := range setA {
		result[key] = struct{}{}
	}

	// Удаляем элементы, присутствующие в setB
	for key := range setB {
		delete(result, key)
	}

	return result
}

func constructFirst(g *grammar.Grammar, first map[string]map[string]struct{}) map[string]map[string]struct{} {
	changed := true

	for changed {
		changed = false

		for nt, rule := range g.Grammar {
			for _, pb := range rule.Rights {
				elem := pb[0]

				newElements := difference(first[elem], first[nt])

				if len(newElements) != 0 {
					first[nt] = union(first[nt], first[elem])

					changed = true
				}
			}
		}
	}

	return first
}

func makeFirst(g *grammar.Grammar) map[string]map[string]struct{} {
	first := make(map[string]map[string]struct{})

	// first(terminal) = terminal
	for _, rule := range g.Grammar {
		first[rule.NonTerminal] = make(map[string]struct{})

		for _, pb := range rule.Rights {
			for _, elem := range pb {
				if grammar.IsTerminal(elem) {
					first[elem] = make(map[string]struct{})
					first[elem][elem] = struct{}{}
				}
			}
		}
	}

	first = constructFirst(g, first)

	// remove terminals from first
	for e := range first {
		if grammar.IsTerminal(e) {
			delete(first, e)
		}
	}

	return first
}

func makeLast(g *grammar.Grammar) map[string]map[string]struct{} {
	return makeFirst(g.Reverse())
}

func checkFollow(
	g *grammar.Grammar,
	follow map[string]map[string]struct{},
	first map[string]map[string]struct{},
) (map[string]map[string]struct{}, bool) {
	changed := false

	isChanged := func(rightRule grammar.ProductionBody) bool {
		for terminal := range first[rightRule[1]] {
			if _, ok := follow[rightRule[0]][terminal]; !ok {
				return true
			}
		}

		return false
	}

	for _, rightRules := range g.Grammar {
		for _, rightRule := range rightRules.Rights {
			if len(rightRule) > 1 {
				changed = isChanged(rightRule)

				follow[rightRule[0]] = union(
					follow[rightRule[0]],
					first[rightRule[1]],
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

func makePrecede(g *grammar.Grammar, last map[string]map[string]struct{}) map[string]map[string]struct{} {
	return makeFollow(g.Reverse(), last)
}

// very bad function, i don't like it
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
			if len(rule) == 2 {
				exists[fmt.Sprintf("%s %s", rule[0], rule[1])] = struct{}{}
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
	terminals := g.ExtractTerminals()

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
	first := makeFirst(g)
	last := makeLast(g)
	follow := makeFollow(g, first)
	precede := makePrecede(g, last)

	matrix := makeBigramMatrix(g, first, last, follow, precede)

	return &Bigramms{
		Matrix: matrix,
		First:  first,
	}
}
