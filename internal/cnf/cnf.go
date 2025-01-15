package cnf

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"github.com/BaldiSlayer/rofl-lab3/pkg/queue"
)

type CNF struct{}

func isNotTerminal(symbols models.SymbolsBtw) bool {
	return !(symbols.S[0] >= 'a' && symbols.S[0] <= 'z')
}

func isTerminal(symbols models.SymbolsBtw) bool {
	return symbols.S[0] >= 'a' && symbols.S[0] <= 'z'
}

func mergeGrammars(parent *grammar.Grammar, child *grammar.Grammar) *grammar.Grammar {
	newGrammar := *parent

	for _, rule := range child.Grammar {
		newGrammar.Grammar[rule.NonTerminal] = rule
	}

	return &newGrammar
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
	nts := make(map[models.SymbolsBtw]struct{}, 2*len(rule.Rights))

	for _, pBody := range rule.Rights {
		for notTerminal := range getNonTerminalsOfProductionBody(pBody) {
			nts[notTerminal] = struct{}{}
		}
	}

	return nts
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

func pbContainsNT(body models.ProductionBody, nt string) bool {
	for _, elem := range body.Body {
		if elem.S == nt {
			return true
		}
	}

	return false
}

func ruleContainsNT(rule models.Rule, nt string) bool {
	for _, right := range rule.Rights {
		if pbContainsNT(right, nt) {
			return true
		}
	}

	return false
}

func deleteRulesWithNT(g *grammar.Grammar, nt string) *grammar.Grammar {
	// удаляем нетерминалы, которые не являются порождающими
	for ngnt := range g.Grammar {
		if nt == ngnt {
			delete(g.Grammar, ngnt)
		}
	}

	// удаляем правила с их вхождением
	for ngnt, rule := range g.Grammar {
		for i, pb := range rule.Rights {
			if pbContainsNT(pb, nt) {
				newRights := append(g.Grammar[ngnt].Rights[:i], g.Grammar[ngnt].Rights[i+1:]...)

				if len(newRights) == 0 {
					delete(g.Grammar, ngnt)

					break
				}

				g.Grammar[ngnt] = models.Rule{
					NonTerminal: g.Grammar[ngnt].NonTerminal,
					Rights:      newRights,
				}
			}
		}
	}

	return g
}

func determineGenerativeness(g *grammar.Grammar) map[string]bool {
	rules := g.GetProductionsSlice()

	isGenerating := make(map[string]bool, len(g.Grammar))
	counter := make([]int, len(rules))
	concernedRules := make(map[string][]int, len(g.Grammar))

	for nt := range g.Grammar {
		isGenerating[nt] = false
	}

	for nt := range g.Grammar {
		a := make([]int, 0)

		for i, rule := range rules {
			if ruleContainsNT(rule, nt) {
				a = append(a, i)
			}
		}

		concernedRules[nt] = a
	}

	for i, rule := range rules {
		counter[i] = len(getNonTerminalsOfRule(rule))
	}

	q := &queue.Queue[string]{}

	for i := 0; i < len(counter); i++ {
		if counter[i] == 0 {
			q.Enqueue(rules[i].NonTerminal)
			isGenerating[rules[i].NonTerminal] = true
		}
	}

	for !q.IsEmpty() {
		cur := q.Dequeue()

		for _, idx := range concernedRules[cur] {
			counter[idx]--

			if counter[idx] == 0 {
				nt := rules[idx].NonTerminal

				isGenerating[nt] = true
				q.Enqueue(nt)
			}
		}
	}

	return isGenerating
}

func deleteNonGenerative(g *grammar.Grammar) *grammar.Grammar {
	isGenerating := determineGenerativeness(g)

	for nt, val := range isGenerating {
		if !val {
			g = deleteRulesWithNT(g, nt)
		}
	}

	return g
}

func findNonReachable(start string, g *grammar.Grammar, visited map[string]struct{}) map[string]struct{} {
	visited[start] = struct{}{}

	for _, rightRule := range g.Grammar[start].Rights {
		for _, smb := range rightRule.Body {
			if _, ok := visited[smb.S]; !ok {
				visited = findNonReachable(smb.S, g, visited)
			}
		}
	}

	return visited
}

func deleteNonReachable(g *grammar.Grammar) *grammar.Grammar {
	visited := make(map[string]struct{})

	visited = findNonReachable(g.Start, g, visited)

	for nt := range g.Grammar {
		if _, ok := visited[nt]; !ok {
			delete(g.Grammar, nt)
		}
	}

	return g
}

func replacePairedTerminals(
	pb models.ProductionBody,
	genNT func() string,
) (models.ProductionBody, []models.Rule) {
	rules := make([]models.Rule, 0)

	if len(pb.Body) == 2 && isTerminal(pb.Body[0]) && isTerminal(pb.Body[1]) {
		if pb.Body[0] == pb.Body[1] {
			name := genNT()

			rules = append(rules, models.Rule{
				NonTerminal: name,
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							pb.Body[0],
						},
					},
				},
			})

			return models.ProductionBody{
				Body: []models.SymbolsBtw{
					{name}, {name},
				},
			}, rules
		}

		// todo разные
	}

	return pb, rules
}

func deletePairedTerminals(g *grammar.Grammar) *grammar.Grammar {
	i := -1

	genNTName := func() string {
		i++
		return fmt.Sprintf("NT_PT_%d", i)
	}

	replacements := make([]models.Rule, 0)

	for nt, rules := range g.Grammar {
		newRights := make([]models.ProductionBody, 0, len(rules.Rights))

		for _, pb := range rules.Rights {
			newPB, r := replacePairedTerminals(pb, genNTName)

			newRights = append(newRights, newPB)
			replacements = append(replacements, r...)
		}

		g.Grammar[nt] = models.Rule{
			NonTerminal: nt,
			Rights:      newRights,
		}
	}

	for _, r := range replacements {
		g.Grammar[r.NonTerminal] = r
	}

	return g
}

func (cnf *CNF) ToCNF(g *grammar.Grammar) *grammar.Grammar {
	transformations := [...]func(*grammar.Grammar) *grammar.Grammar{
		deleteLongRules,
		deleteChainRules,
		deleteNonGenerative,
		deleteNonReachable,
		deletePairedTerminals,
	}

	// TODO it looks bad, I don't like it, but writing 7 function calls and declaring
	// variables for them was even more annoying for me.
	for _, transformation := range transformations {
		g = transformation(g)
	}

	return g
}
