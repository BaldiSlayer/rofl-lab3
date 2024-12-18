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

type Grammar struct {
}

func New() *Grammar {
	return &Grammar{}
}

func removeSpacesAndStrip(s string) string {
	withoutSpaces := strings.ReplaceAll(s, " ", "")
	trimmed := strings.TrimFunc(withoutSpaces, unicode.IsSpace)

	return trimmed
}

func parseT(_ <-chan byte, prevSymbol byte) SymbolsBtw {
	return SymbolsBtw{
		s: string(prevSymbol),
	}
}

func parseBetweenBrackets(in <-chan byte, _ byte) []SymbolsBtw {
	var sb strings.Builder

	sb.WriteByte('[')

	for c := range in {
		sb.WriteByte(c)

		if c == ']' {
			return []SymbolsBtw{
				{
					s: sb.String(),
				},
			}
		}
	}

	return nil
}

func parseCapitalDigit(in <-chan byte, prevSymbol byte) []SymbolsBtw {
	res := <-in

	if res >= '0' && res <= '9' {
		return []SymbolsBtw{
			{
				s: string(prevSymbol) + string(res),
			},
		}
	}

	return []SymbolsBtw{
		{
			s: string(prevSymbol),
		},
		{
			s: string(res),
		},
	}
}

func parseNT(in <-chan byte, prevSymbol byte) []SymbolsBtw {
	if prevSymbol == '[' {
		return parseBetweenBrackets(in, prevSymbol)
	}

	return parseCapitalDigit(in, prevSymbol)
}

func parse(in <-chan byte, out chan<- SymbolsBtw) {
	for c := range in {
		if c >= 'a' && c <= 'z' {
			out <- parseT(in, c)

			continue
		}

		for _, nt := range parseNT(in, c) {
			out <- nt
		}
	}

	close(out)
}

func parseProductionBody(s string) ProductionBody {
	in := make(chan byte, 0)
	out := make(chan SymbolsBtw, 0)

	go parse(in, out)

	go func() {
		for _, hui := range s {
			in <- byte(hui)
		}

		close(in)
	}()

	body := make([]SymbolsBtw, 0)

	for c := range out {
		body = append(body, c)
	}

	return ProductionBody{
		body: body,
	}
}

func parseRight(s string) []ProductionBody {
	pbs := make([]ProductionBody, 0)

	for _, production := range strings.Split(s, "|") {
		trimmed := strings.TrimFunc(production, unicode.IsSpace)

		pbs = append(pbs, parseProductionBody(trimmed))
	}

	return pbs
}

func parseLine(s string) Rule {
	s = removeSpacesAndStrip(s)

	split := strings.Split(s, "->")

	parseRight(split[1])

	return Rule{
		nonTerminal: split[0],
	}
}

func (g *Grammar) ParseLines(lines []string) {
	for _, line := range lines {
		parseLine(line)
	}
}
