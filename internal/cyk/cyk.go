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

func calcDP(d map[string][][]bool, rightRules models.Rule, i, j int) bool {
	for _, rightRule := range rightRules.Rights {
		for k := i + 1; k < j; k++ {
			if d[rightRule.Body[0].S][i][k] && d[rightRule.Body[1].S][k][j] {
				return true
			}
		}
	}

	return false
}

func (c *CYK) Check(word string) bool {
	dp := make(map[string][][]bool)

	for _, rule := range c.g.Grammar {
		dp[rule.NonTerminal] = make([][]bool, len(word)+1)

		for i := range dp[rule.NonTerminal] {
			dp[rule.NonTerminal][i] = make([]bool, len(word)+1)
		}
	}

	for i := 0; i < len(word); i++ {
		for _, rightRules := range c.terminalRules {
			for _, rightRule := range rightRules.Rights {
				dp[rightRules.NonTerminal][i][i+1] = isOneTermRule(rightRule, word[i])
			}
		}
	}

	for m := 2; m < len(word)+1; m++ {
		for i := 0; i < len(word)-m+1; i++ {
			j := i + m

			for _, rightRules := range c.nonTerminalRules {
				dp[rightRules.NonTerminal][i][j] = dp[rightRules.NonTerminal][i][j] || calcDP(dp, rightRules, i, j)
			}
		}
	}

	return dp[c.startingTerm][0][len(word)]
}
