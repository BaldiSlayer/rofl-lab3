package parser

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"strings"
	"unicode"
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

func (p *Parser) parseBetweenBrackets(s string) models.SymbolsBtw {
	var sb strings.Builder

	sb.WriteByte('[')

	for _, sym := range s {
		sb.WriteByte(byte(sym))

		if sym == ']' {
			p.pos += sb.Len() + 1

			return models.SymbolsBtw{
				S: sb.String(),
			}
		}
	}

	return models.SymbolsBtw{}
}

func isNumeric(symbol byte) bool {
	return symbol >= '0' && symbol <= '9'
}

func (p *Parser) parseCapitals(s string) models.SymbolsBtw {
	if len(s) > 1 && isNumeric(s[1]) {
		p.pos++

		return models.SymbolsBtw{
			S: string(s[0]) + string(s[p.pos+1]),
		}
	}

	return models.SymbolsBtw{
		S: string(s[0]),
	}
}

func (p *Parser) parseProductionBody(s string) models.ProductionBody {
	body := make([]models.SymbolsBtw, 0, len(s))

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
		body = append(body, models.SymbolsBtw{
			S: string(s[i]),
		})

	}

	return models.ProductionBody{
		Body: body,
	}
}

func (p *Parser) parseRight(s string) []models.ProductionBody {
	pbs := make([]models.ProductionBody, 0)

	for _, production := range strings.Split(s, "|") {
		trimmed := strings.TrimFunc(production, unicode.IsSpace)

		pbs = append(pbs, p.parseProductionBody(trimmed))
	}

	return pbs
}

func (p *Parser) parseLine(s string) models.Rule {
	s = removeSpacesAndStrip(s)

	split := strings.Split(s, "->")

	return models.Rule{
		NonTerminal: split[0],
		Rights:      p.parseRight(split[1]),
	}
}

func (p *Parser) parseLines(lines []string) []models.Rule {
	rules := make([]models.Rule, 0, len(lines))

	for _, line := range lines {
		rules = append(rules, p.parseLine(line))
	}

	return rules
}

func (p *Parser) Parse(input string) *grammar.Grammar {
	lines := strings.Split(input, "\n")

	rules := p.parseLines(lines)

	return grammar.New(rules)
}