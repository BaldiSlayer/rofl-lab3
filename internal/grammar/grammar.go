package grammar

import "github.com/BaldiSlayer/rofl-lab3/internal/models"

type Grammar struct {
	Start   string
	Grammar map[string]models.Rule
}

func isTerminal(symbols string) bool {
	return symbols[0] >= 'a' && symbols[0] <= 'z'
}

func extractTerminalsFromRule(rule models.ProductionBody) []string {
	terminals := make([]string, 0)

	for _, smb := range rule {
		if isTerminal(smb) {
			terminals = append(terminals, smb)
		}
	}

	return terminals
}

func (g *Grammar) ExtractTerminals() []string {
	terminals := make([]string, 0, len(g.Grammar))

	for _, rules := range g.Grammar {
		for _, rule := range rules.Rights {
			terminals = append(terminals, extractTerminalsFromRule(rule)...)
		}
	}

	return terminals
}

// GetProductionsSlice creates a slide with rules. Necessary for a constant order
func (g *Grammar) GetProductionsSlice() []models.Rule {
	rules := make([]models.Rule, 0, len(g.Grammar))

	for _, pbs := range g.Grammar {
		for _, rightRule := range pbs.Rights {
			rules = append(rules, models.Rule{
				NonTerminal: pbs.NonTerminal,
				Rights: []models.ProductionBody{
					rightRule,
				},
			})
		}
	}

	return rules
}

func New(rules []models.Rule) *Grammar {
	g := make(map[string]models.Rule, len(rules))

	for _, rule := range rules {
		if _, ok := g[rule.NonTerminal]; !ok {
			g[rule.NonTerminal] = rule

			continue
		}

		g[rule.NonTerminal] = models.Rule{
			NonTerminal: rule.NonTerminal,
			Rights:      append(g[rule.NonTerminal].Rights, rule.Rights...),
		}
	}

	// todo убрать хардкод стартового, вынести это в параметры
	return &Grammar{
		Grammar: g,
		Start:   "S",
	}
}
