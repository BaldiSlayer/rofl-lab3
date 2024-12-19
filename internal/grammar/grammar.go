package grammar

import "github.com/BaldiSlayer/rofl-lab3/internal/models"

type Grammar struct {
	start   string
	grammar map[string]models.Rule
}

func New(rules []models.Rule) *Grammar {
	g := make(map[string]models.Rule, len(rules))

	for _, rule := range rules {
		g[rule.NonTerminal] = rule
	}

	return &Grammar{
		grammar: g,
	}
}
