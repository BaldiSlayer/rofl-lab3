package cyk

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
)

type CYK struct {
	g                *grammar.Grammar
	terminalRules    []models.Rule
	nonTerminalRules []models.Rule
	startingTerm     string
}

func New(g *grammar.Grammar) *CYK {
	terminalRules := make([]models.Rule, 0, len(g.Grammar))
	nonTerminalRules := make([]models.Rule, 0, len(g.Grammar))

	for _, rightRules := range g.Grammar {
		for _, rightRule := range rightRules.Rights {
			rule := models.Rule{
				NonTerminal: rightRules.NonTerminal,
				Rights: []models.ProductionBody{
					rightRule,
				},
			}

			if len(rightRule.Body) == 1 {
				terminalRules = append(terminalRules, rule)

				continue
			}

			nonTerminalRules = append(nonTerminalRules, rule)
		}
	}

	return &CYK{
		g:                g,
		terminalRules:    terminalRules,
		nonTerminalRules: nonTerminalRules,
		startingTerm:     g.Start,
	}
}

func isOneTermRule(rule models.ProductionBody, c uint8) bool {
	return len(rule.Body) == 1 && string(c) == rule.Body[0].S
}

func dp(d map[string][][]bool, rightRules models.Rule, i, j int) bool {
	for _, rightRule := range rightRules.Rights {
		for k := i; k < j; k++ {
			if d[rightRule.Body[0].S][i][k] && d[rightRule.Body[1].S][k+1][j] {
				return true
			}
		}
	}

	return false
}

func (c *CYK) Check(word string) bool {
	d := make(map[string][][]bool)

	for _, rule := range c.g.Grammar {
		d[rule.NonTerminal] = make([][]bool, len(word))

		for i := range d[rule.NonTerminal] {
			d[rule.NonTerminal][i] = make([]bool, len(word))
		}
	}

	for i := 0; i < len(word); i++ {
		for _, rightRules := range c.terminalRules {
			for _, rightRule := range rightRules.Rights {
				d[rightRules.NonTerminal][i][i] = isOneTermRule(rightRule, word[i])
			}
		}
	}

	for m := 1; m < len(word); m++ {
		for i := 0; i < len(word)-m; i++ {
			j := i + m

			if i == 3 && m == 2 {
				u := 0
				_ = u
			}

			for _, rightRules := range c.nonTerminalRules {
				d[rightRules.NonTerminal][i][j] = dp(d, rightRules, i, j)
			}
		}
	}

	return d[c.startingTerm][0][len(word)-1]
}
