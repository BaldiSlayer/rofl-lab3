package cnf

import (
	"fmt"

	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"

	"github.com/BaldiSlayer/rofl-lab3/pkg/queue"
)

type CNF struct{}

func mergeGrammars(parent *grammar.Grammar, child *grammar.Grammar) *grammar.Grammar {
	for _, rule := range child.Grammar {
		parent.Grammar[rule.NonTerminal] = rule
	}

	return parent
}

func newIDGetter() func() int {
	// чтобы начать с 0
	i := -1

	return func() int {
		i++
		return i
	}
}

func deleteLongPart(rule grammar.Rule, idGetter func() int) []grammar.Rule {
	rules := make([]grammar.Rule, 0)

	for _, r := range rule.Rights {
		if len(r) <= 2 {
			rules = append(rules, grammar.Rule{
				NonTerminal: rule.NonTerminal,
				Rights:      []grammar.ProductionBody{r},
			})

			continue
		}

		newNT := fmt.Sprintf("[new_NT_%d]", idGetter())

		shortRule := grammar.Rule{
			NonTerminal: rule.NonTerminal,
			Rights: []grammar.ProductionBody{
				{
					r[0],
					newNT,
				},
			},
		}

		newRule := grammar.Rule{
			NonTerminal: newNT,
			Rights:      []grammar.ProductionBody{r[1:]},
		}

		rules = append(rules, shortRule)
		rules = append(rules, deleteLongPart(newRule, idGetter)...)
	}

	return rules
}

func deleteLongRules(g *grammar.Grammar) *grammar.Grammar {
	rules := make([]grammar.Rule, 0)

	idGetter := newIDGetter()

	for _, rule := range g.Grammar {
		rules = append(rules, deleteLongPart(rule, idGetter)...)
	}

	return grammar.New(rules, g.Start)
}

func getNonTerminalsOfProductionBody(pBody grammar.ProductionBody) map[string]struct{} {
	nts := make(map[string]struct{}, 0)

	for _, symbol := range pBody {
		if grammar.IsNotTerminal(symbol) {
			nts[symbol] = struct{}{}
		}
	}

	return nts
}

func getNonTerminalsOfRule(rule grammar.Rule) map[string]struct{} {
	nts := make(map[string]struct{}, 2*len(rule.Rights))

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
		Grammar: make(map[string]grammar.Rule),
	}

	visited[nt] = struct{}{}

	for symbol := range getNonTerminalsOfRule(g.Grammar[nt]) {
		if _, ok := visited[symbol]; !ok {
			newGrammar = mergeGrammars(
				newGrammar,
				deleteChainRulesIteratively(symbol, g, visited),
			)
		}
	}

	newRule := grammar.Rule{
		NonTerminal: nt,
	}

	for _, pBody := range g.Grammar[nt].Rights {
		// если тело продукции - цепное правило, то все его правила
		// прикрепляем к нетерминалу nt
		if len(pBody) == 1 && grammar.IsNotTerminal(pBody[0]) {
			ntPBs := make([]grammar.ProductionBody, 0)

			if _, ok := newGrammar.Grammar[pBody[0]]; ok {
				ntPBs = newGrammar.Grammar[pBody[0]].Rights
			} else {
				ntPBs = g.Grammar[pBody[0]].Rights
			}

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

func pbContainsNT(body grammar.ProductionBody, nt string) bool {
	for _, elem := range body {
		if elem == nt {
			return true
		}
	}

	return false
}

func ruleContainsNT(rule grammar.Rule, nt string) bool {
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

				g.Grammar[ngnt] = grammar.Rule{
					NonTerminal: g.Grammar[ngnt].NonTerminal,
					Rights:      newRights,
				}
			}
		}
	}

	return g
}

func determineGenerativeness(g *grammar.Grammar) map[string]bool {
	rules := g.GetRulesSlice()

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
		for _, smb := range rightRule {
			if _, ok := visited[smb]; !ok {
				visited = findNonReachable(smb, g, visited)
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
	pb grammar.ProductionBody,
	genNT func(terminal string) string,
) (grammar.ProductionBody, []grammar.Rule) {
	rules := make([]grammar.Rule, 0)
	pBody := make(grammar.ProductionBody, 0)

	checkSmb := func(smb string) {
		name := smb

		if grammar.IsTerminal(smb) {
			if len(rules) != 0 && pb[0] == pb[1] {
				pBody = append(pBody, rules[0].NonTerminal)

				return
			}

			name = genNT(smb)

			rules = append(rules, grammar.Rule{
				NonTerminal: name,
				Rights: []grammar.ProductionBody{
					{
						smb,
					},
				},
			})
		}

		pBody = append(pBody, name)
	}

	if len(pb) == 2 {
		checkSmb(pb[0])
		checkSmb(pb[1])

		return pBody, rules
	}

	return pb, rules
}

func deletePairedTerminals(g *grammar.Grammar) *grammar.Grammar {
	genNTName := func(terminal string) string {
		return fmt.Sprintf("[NT_PT_%s]", terminal)
	}

	replacements := make([]grammar.Rule, 0)

	for nt, rules := range g.Grammar {
		newRights := make([]grammar.ProductionBody, 0, len(rules.Rights))

		for _, pb := range rules.Rights {
			newPB, r := replacePairedTerminals(pb, genNTName)

			newRights = append(newRights, newPB)
			replacements = append(replacements, r...)
		}

		g.Grammar[nt] = grammar.Rule{
			NonTerminal: nt,
			Rights:      newRights,
		}
	}

	for _, r := range replacements {
		g.Grammar[r.NonTerminal] = r
	}

	return g
}

func replaceAloneTerminals(
	pb grammar.ProductionBody,
	genNT func(terminal string) string,
) (grammar.ProductionBody, []grammar.Rule) {
	rules := make([]grammar.Rule, 0)
	pBody := make(grammar.ProductionBody, 0)

	if len(pb) == 1 && grammar.IsTerminal(pb[0]) {
		name := genNT(pb[0])

		pBody = append(pBody, name)

		rules = append(rules, grammar.Rule{
			NonTerminal: name,
			Rights: []grammar.ProductionBody{
				{
					pb[0],
				},
			},
		})

		return pBody, rules
	}

	return pb, rules
}

func deleteAloneTerminals(g *grammar.Grammar) *grammar.Grammar {
	genNTName := func(terminal string) string {
		return fmt.Sprintf("[NT_ALONE_%s]", terminal)
	}

	replacements := make([]grammar.Rule, 0)

	for nt, rules := range g.Grammar {
		if len(rules.Rights) > 1 {
			newRights := make([]grammar.ProductionBody, 0, len(rules.Rights))

			for _, pb := range rules.Rights {
				newPB, r := replaceAloneTerminals(pb, genNTName)

				newRights = append(newRights, newPB)
				replacements = append(replacements, r...)
			}

			g.Grammar[nt] = grammar.Rule{
				NonTerminal: nt,
				Rights:      newRights,
			}
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
		deleteAloneTerminals,
	}

	// TODO it looks bad, I don't like it, but writing 5 function calls and declaring
	// variables for them was even more annoying for me.
	for _, transformation := range transformations {
		g = transformation(g)
	}

	return g
}
