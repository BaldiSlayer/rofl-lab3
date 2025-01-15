package bigramms

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	dest := map[string]struct{}{"a": {}, "b": {}}
	src := map[string]struct{}{"b": {}, "c": {}}
	expected := map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithEmptySrc(t *testing.T) {
	dest := map[string]struct{}{"x": {}, "y": {}}
	src := map[string]struct{}{}
	expected := map[string]struct{}{
		"x": {},
		"y": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithEmptyDest(t *testing.T) {
	dest := map[string]struct{}{}
	src := map[string]struct{}{"m": {}, "n": {}}
	expected := map[string]struct{}{
		"m": {},
		"n": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithIdenticalMaps(t *testing.T) {
	dest := map[string]struct{}{"p": {}}
	src := map[string]struct{}{"p": {}}
	expected := map[string]struct{}{
		"p": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestBigramms_Build(t *testing.T) {
	g := grammar.Grammar{
		Start: "S",
		Grammar: map[string]models.Rule{
			"S": {
				NonTerminal: "S",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"A",
							},
							{
								"A",
							},
						},
					},
				},
			},
			"A": {
				NonTerminal: "A",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"a",
							},
						},
					},
				},
			},
		},
	}

	b := &Bigramms{}

	b.Build(&g)

	u := 0
	_ = u
}
