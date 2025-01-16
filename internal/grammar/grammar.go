package grammar

import (
	"strings"
)

type ProductionBody []string

type Rule struct {
	NonTerminal string
	Rights      []ProductionBody
}

type Grammar struct {
	Start   string
	Grammar map[string]Rule
}

func isTerminal(symbols string) bool {
	return symbols[0] >= 'a' && symbols[0] <= 'z'
}

func extractTerminalsFromRule(rule ProductionBody) []string {
	terminals := make([]string, 0)

	for _, smb := range rule {
		if isTerminal(smb) {
			terminals = append(terminals, smb)
		}
	}

	return terminals
}

func uniqify(input []string) []string {
	uniqueMap := make(map[string]struct{}, len(input))
	uniqueSlice := make([]string, 0, len(input))

	for _, item := range input {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = struct{}{}
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice
}

func (g *Grammar) ExtractTerminals() []string {
	terminals := make([]string, 0, len(g.Grammar))

	for _, rules := range g.Grammar {
		for _, rule := range rules.Rights {
			terminals = append(terminals, extractTerminalsFromRule(rule)...)
		}
	}

	return uniqify(terminals)
}

// GetRulesSlice creates a slide with rules. Necessary for a constant order
func (g *Grammar) GetRulesSlice() []Rule {
	rules := make([]Rule, 0, len(g.Grammar))

	for _, pbs := range g.Grammar {
		for _, rightRule := range pbs.Rights {
			rules = append(rules, Rule{
				NonTerminal: pbs.NonTerminal,
				Rights: []ProductionBody{
					rightRule,
				},
			})
		}
	}

	return rules
}

func (g *Grammar) String() string {
	var sb strings.Builder

	for _, rule := range g.Grammar {
		sb.WriteString(rule.NonTerminal)
		sb.WriteString(" -> ")

		for i, pb := range rule.Rights {
			smth := ""

			for _, smb := range pb {
				smth += smb
			}

			sb.WriteString(smth)

			if i < len(rule.Rights)-1 {
				sb.WriteString(" | ")
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func New(rules []Rule, startSymbol string) *Grammar {
	g := make(map[string]Rule, len(rules))

	for _, rule := range rules {
		if _, ok := g[rule.NonTerminal]; !ok {
			g[rule.NonTerminal] = rule

			continue
		}

		g[rule.NonTerminal] = Rule{
			NonTerminal: rule.NonTerminal,
			Rights:      append(g[rule.NonTerminal].Rights, rule.Rights...),
		}
	}

	return &Grammar{
		Grammar: g,
		Start:   startSymbol,
	}
}
