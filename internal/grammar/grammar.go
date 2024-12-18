package grammar

import (
	"strings"
	"unicode"
)

type SymbolsBtw struct {
	s string
}

type ProductionBody struct {
	body []SymbolsBtw
}

type Rule struct {
	nonTerminal string
	rights      []ProductionBody
}

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

func (p *Parser) parseBetweenBrackets(s string) SymbolsBtw {
	var sb strings.Builder

	sb.WriteByte('[')

	for _, sym := range s {
		sb.WriteByte(byte(sym))

		if sym == ']' {
			p.pos += sb.Len() + 1

			return SymbolsBtw{
				s: sb.String(),
			}
		}
	}

	return SymbolsBtw{}
}

func isNumeric(symbol byte) bool {
	return symbol >= '0' && symbol <= '9'
}

func (p *Parser) parseCapitals(s string) SymbolsBtw {
	if len(s) > 1 && isNumeric(s[1]) {
		p.pos++

		return SymbolsBtw{
			s: string(s[0]) + string(s[p.pos+1]),
		}
	}

	return SymbolsBtw{
		s: string(s[0]),
	}
}

func (p *Parser) parseProductionBody(s string) ProductionBody {
	body := make([]SymbolsBtw, 0, len(s))

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
		body = append(body, SymbolsBtw{
			s: string(s[i]),
		})

	}

	return ProductionBody{
		body: body,
	}
}

func (p *Parser) parseRight(s string) []ProductionBody {
	pbs := make([]ProductionBody, 0)

	for _, production := range strings.Split(s, "|") {
		trimmed := strings.TrimFunc(production, unicode.IsSpace)

		pbs = append(pbs, p.parseProductionBody(trimmed))
	}

	return pbs
}

func (p *Parser) parseLine(s string) Rule {
	s = removeSpacesAndStrip(s)

	split := strings.Split(s, "->")

	return Rule{
		nonTerminal: split[0],
		rights:      p.parseRight(split[1]),
	}
}

func (p *Parser) ParseLines(lines []string) {
	for _, line := range lines {
		p.parseLine(line)
	}
}
