package models

type ProductionBody struct {
	Body []string
}

type Rule struct {
	NonTerminal string
	Rights      []ProductionBody
}
