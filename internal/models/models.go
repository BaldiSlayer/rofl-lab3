package models

type ProductionBody []string

type Rule struct {
	NonTerminal string
	Rights      []ProductionBody
}
