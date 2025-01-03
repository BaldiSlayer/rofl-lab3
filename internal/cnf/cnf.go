package cnf

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
)

type CNF struct {
}

func newIDGetter() func() int {
	// чтобы начать с 0
	i := -1

	return func() int {
		i++
		return i
	}
}

func deleteLongPart(rule models.Rule, idGetter func() int) []models.Rule {
	rules := make([]models.Rule, 0)

	for _, r := range rule.Rights {
		if len(r.Body) <= 2 {
			rules = append(rules, models.Rule{
				NonTerminal: rule.NonTerminal,
				Rights:      []models.ProductionBody{r},
			})

			continue
		}

		newNT := fmt.Sprintf("[new_NT_%d]", idGetter())

		shortRule := models.Rule{
			NonTerminal: rule.NonTerminal,
			Rights: []models.ProductionBody{
				{
					Body: []models.SymbolsBtw{
						{
							r.Body[0].S,
						},
						{
							newNT,
						},
					},
				},
			},
		}

		newRule := models.Rule{
			NonTerminal: newNT,
			Rights: []models.ProductionBody{
				{
					Body: r.Body[1:],
				},
			},
		}

		rules = append(rules, shortRule)
		rules = append(rules, deleteLongPart(newRule, idGetter)...)
	}

	return rules
}

func deleteLongRules(g *grammar.Grammar) *grammar.Grammar {
	rules := make([]models.Rule, 0)

	idGetter := newIDGetter()

	for _, rule := range g.Grammar {
		rules = append(rules, deleteLongPart(rule, idGetter)...)
	}

	return grammar.New(rules)
}

func isNotTerminal(symbols models.SymbolsBtw) bool {
	return !(symbols.S[0] >= 'a' && symbols.S[0] <= 'z')
}

func getNonTerminalsOfProductionBody(pBody models.ProductionBody) map[models.SymbolsBtw]struct{} {
	nts := make(map[models.SymbolsBtw]struct{}, 0)

	for _, symbol := range pBody.Body {
		if isNotTerminal(symbol) {
			nts[symbol] = struct{}{}
		}
	}

	return nts
}

func getNonTerminalsOfRule(rule models.Rule) map[models.SymbolsBtw]struct{} {
	nts := make(map[models.SymbolsBtw]struct{}, 0)

	for _, pBody := range rule.Rights {
		for notTerminal := range getNonTerminalsOfProductionBody(pBody) {
			nts[notTerminal] = struct{}{}
		}
	}

	return nts
}

func mergeGrammars(parent *grammar.Grammar, child *grammar.Grammar) *grammar.Grammar {
	newGrammar := *parent

	for _, rule := range child.Grammar {
		newGrammar.Grammar[rule.NonTerminal] = rule
	}

	return &newGrammar
}

func deleteChainRulesIteratively(nt string, g *grammar.Grammar, visited map[string]struct{}) *grammar.Grammar {
	newGrammar := &grammar.Grammar{
		Start:   nt,
		Grammar: make(map[string]models.Rule),
	}

	visited[nt] = struct{}{}

	for symbol := range getNonTerminalsOfRule(g.Grammar[nt]) {
		if _, ok := visited[symbol.S]; !ok {
			newGrammar = mergeGrammars(
				newGrammar,
				deleteChainRulesIteratively(symbol.S, g, visited),
			)
		}
	}

	newRule := models.Rule{
		NonTerminal: nt,
	}

	for _, pBody := range g.Grammar[nt].Rights {
		// если тело продукции - цепное правило, то все его правила
		// прикрепляем к нетерминалу nt
		if len(pBody.Body) == 1 && isNotTerminal(pBody.Body[0]) {
			ntPBs := g.Grammar[pBody.Body[0].S].Rights
			newRule.Rights = append(newRule.Rights, ntPBs...)

			continue
		}

		newRule.Rights = append(newRule.Rights, pBody)
	}

	newGrammar.Grammar[nt] = newRule

	return newGrammar
}

func deleteChainRules(g *grammar.Grammar) *grammar.Grammar {
	visited := make(map[string]struct{}, len(g.Grammar))

	newGrammar := deleteChainRulesIteratively(g.Start, g, visited)

	return newGrammar
}

func (cnf *CNF) ToCNF(g *grammar.Grammar) *grammar.Grammar {
	transformations := []func(*grammar.Grammar) *grammar.Grammar{
		deleteLongRules,
		deleteChainRules,
	}

	// TODO it looks bad, I don't like it, but writing 7 function calls and declaring
	// variables for them was even more annoying for me.
	for _, transformation := range transformations {
		g = transformation(g)
	}

	return g
}
