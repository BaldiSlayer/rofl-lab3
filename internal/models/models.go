package models

type SymbolsBtw struct {
	S string
}

type ProductionBody struct {
	Body []SymbolsBtw
}

type Rule struct {
	NonTerminal string
	Rights      []ProductionBody
}
