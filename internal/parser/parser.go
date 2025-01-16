package parser

import (
	"strings"
	"unicode"

	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
)

type Parser struct {
	pos int
}

func New() *Parser {
	return &Parser{}
}

func removeSpacesAndStrip(s string) string {
	withoutSpaces := strings.ReplaceAll(s, " ", "")
	trimmed := strings.TrimFunc(withoutSpaces, unicode.IsSpace)

	return trimmed
}

func (p *Parser) parseBetweenBrackets(s string) string {
	var sb strings.Builder

	for _, sym := range s {
		sb.WriteByte(byte(sym))

		if sym == ']' {
			p.pos += sb.Len() - 1

			return sb.String()
		}
	}

	return ""
}

func isNumeric(symbol byte) bool {
	return symbol >= '0' && symbol <= '9'
}

func (p *Parser) parseCapitals(s string) string {
	if len(s) > 1 && isNumeric(s[1]) {
		p.pos++

		return string(s[0]) + string(s[1])
	}

	return string(s[0])
}

func (p *Parser) parseProductionBody(s string) grammar.ProductionBody {
	body := make([]string, 0, len(s))

	p.pos = 0

	for ; p.pos < len(s); p.pos++ {
		i := p.pos

		if s[i] == '[' {
			body = append(body, p.parseBetweenBrackets(s[i:]))

			continue
		}

		if s[i] >= 'A' && s[i] <= 'Z' {
			body = append(body, p.parseCapitals(s[i:]))

			continue
		}

		// I will panic on wrong data, should I check it?
		body = append(body, string(s[i]))
	}

	return body
}

func (p *Parser) parseRight(s string) []grammar.ProductionBody {
	pbs := make([]grammar.ProductionBody, 0)

	for _, production := range strings.Split(s, "|") {
		trimmed := strings.TrimFunc(production, unicode.IsSpace)

		pbs = append(pbs, p.parseProductionBody(trimmed))
	}

	return pbs
}

func (p *Parser) parseLine(s string) grammar.Rule {
	s = removeSpacesAndStrip(s)

	split := strings.Split(s, "->")

	return grammar.Rule{
		NonTerminal: split[0],
		Rights:      p.parseRight(split[1]),
	}
}

func (p *Parser) parseLines(lines []string) []grammar.Rule {
	rules := make([]grammar.Rule, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		rules = append(rules, p.parseLine(line))
	}

	return rules
}

func (p *Parser) Parse(input string, startSymbol string) *grammar.Grammar {
	lines := strings.Split(input, "\n")

	rules := p.parseLines(lines)

	return grammar.New(rules, startSymbol)
}
