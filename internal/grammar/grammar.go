package grammar

import "github.com/BaldiSlayer/rofl-lab3/internal/models"

type Grammar struct {
	Start   string
	Grammar map[string]models.Rule
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

	return &Grammar{
		Grammar: g,
	}
}
