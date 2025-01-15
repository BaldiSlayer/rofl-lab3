package cyk

import (
	"testing"

	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCYK_Check_1(t *testing.T) {
	g := grammar.Grammar{
		Start: "S",
		Grammar: map[string]models.Rule{
			"S": {
				NonTerminal: "S",
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

	type fields struct {
		g *grammar.Grammar
	}

	tests := []struct {
		name   string
		fields fields
		args   string
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				g: &g,
			},
			args: "a",
			want: true,
		},
		{
			name: "2",
			fields: fields{
				g: &g,
			},
			args: "aa",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.g)

			res := c.Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}

func TestCYK_Check_2(t *testing.T) {
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

	type fields struct {
		g *grammar.Grammar
	}

	tests := []struct {
		name   string
		fields fields
		args   string
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				g: &g,
			},
			args: "a",
			want: false,
		},
		{
			name: "2",
			fields: fields{
				g: &g,
			},
			args: "aa",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.g)

			res := c.Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}

func TestCYK_Check_PSP(t *testing.T) {
	g := grammar.Grammar{
		Start: "S",
		Grammar: map[string]models.Rule{
			"S": {
				NonTerminal: "S",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"B",
							},
							{
								"B",
							},
						},
					},
					{
						Body: []models.SymbolsBtw{
							{
								"C",
							},
							{
								"D",
							},
						},
					},
				},
			},
			"B": {
				NonTerminal: "B",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"B",
							},
							{
								"B",
							},
						},
					},
					{
						Body: []models.SymbolsBtw{
							{
								"C",
							},
							{
								"D",
							},
						},
					},
				},
			},
			"C": {
				NonTerminal: "C",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"(",
							},
						},
					},
				},
			},
			"D": {
				NonTerminal: "D",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"B",
							},
							{
								"E",
							},
						},
					},
					{
						Body: []models.SymbolsBtw{
							{
								")",
							},
						},
					},
				},
			},
			"E": {
				NonTerminal: "E",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								")",
							},
						},
					},
				},
			},
		},
	}

	type fields struct {
		g *grammar.Grammar
	}

	tests := []struct {
		name   string
		fields fields
		args   string
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				g: &g,
			},
			args: "a",
			want: false,
		},
		{
			name: "2",
			fields: fields{
				g: &g,
			},
			args: "aa",
			want: false,
		},
		{
			name: "3",
			fields: fields{
				g: &g,
			},
			args: "())",
			want: false,
		},
		{
			name: "3",
			fields: fields{
				g: &g,
			},
			args: "()",
			want: true,
		},
		{
			name: "4",
			fields: fields{
				g: &g,
			},
			args: "(())",
			want: true,
		},
		{
			name: "5",
			fields: fields{
				g: &g,
			},
			args: "()()()()()()()()()()()()()(()()()()()()()()",
			want: false,
		},
		{
			name: "6",
			fields: fields{
				g: &g,
			},
			args: "()(())",
			want: true,
		},
		{
			name: "7",
			fields: fields{
				g: &g,
			},
			args: "(())()",
			want: true,
		},
		{
			name: "8",
			fields: fields{
				g: &g,
			},
			args: "()(())()",
			want: true,
		},
		{
			name: "9",
			fields: fields{
				g: &g,
			},
			args: "()())",
			want: false,
		},
		{
			name: "10",
			fields: fields{
				g: &g,
			},
			args: "()))",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.g)

			res := c.Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}
