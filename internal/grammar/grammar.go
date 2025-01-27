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

// unifySlice returns slice without repetitions
func unifySlice(input []string) []string {
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

func reverseStringSlice(slice []string) []string {
	reversed := make([]string, len(slice))

	for i, s := range slice {
		reversed[len(slice)-1-i] = s
	}

	return reversed
}

func IsTerminal(symbols string) bool {
	return symbols[0] >= 'a' && symbols[0] <= 'z'
}

func IsNotTerminal(symbols string) bool {
	return !IsTerminal(symbols)
}

// extractTerminalsFromPB extracts terminals from production body
func extractTerminalsFromPB(pb ProductionBody) []string {
	terminals := make([]string, 0)

	for _, smb := range pb {
		if IsTerminal(smb) {
			terminals = append(terminals, smb)
		}
	}

	return terminals
}

// ExtractTerminals extracts grammar terminals
func (g *Grammar) ExtractTerminals() []string {
	terminals := make([]string, 0, len(g.Grammar))

	for _, rules := range g.Grammar {
		for _, rule := range rules.Rights {
			terminals = append(terminals, extractTerminalsFromPB(rule)...)
		}
	}

	return unifySlice(terminals)
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

// Reverse returns grammar with revered production bodies
func (g *Grammar) Reverse() *Grammar {
	var gram Grammar

	gram.Grammar = make(map[string]Rule)

	for nt, rules := range g.Grammar {
		rights := make([]ProductionBody, 0)

		for _, rule := range rules.Rights {
			rights = append(rights, reverseStringSlice(rule))
		}

		gram.Grammar[nt] = Rule{
			NonTerminal: nt,
			Rights:      rights,
		}
	}

	return &gram
}

// String translates grammar to string
func (g *Grammar) String() string {
	var sb strings.Builder

	for _, rule := range g.Grammar {
		sb.WriteString(rule.NonTerminal)
		sb.WriteString(" -> ")

		for i, pb := range rule.Rights {
			pbStr := ""

			for _, smb := range pb {
				pbStr += smb
			}

			sb.WriteString(pbStr)

			if i < len(rule.Rights)-1 {
				sb.WriteString(" | ")
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// New creates grammar object from rules slice and startSymbol
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
