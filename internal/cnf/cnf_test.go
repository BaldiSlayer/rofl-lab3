package cnf

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
	"testing"

	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"github.com/stretchr/testify/require"
)

func Test_deleteLongPart(t *testing.T) {
	type args struct {
		rule     models.Rule
		idGetter func() int
	}
	tests := []struct {
		name string
		args args
		want []models.Rule
	}{
		{
			name: "1",
			args: args{
				rule: models.Rule{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"S",
							"S",
							"a",
						},
					},
				},
				idGetter: newIDGetter(),
			},
			want: []models.Rule{
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"S",
							"[new_NT_0]",
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							"S",
							"a",
						},
					},
				},
			},
		},
		{
			name: "2",
			args: args{
				rule: models.Rule{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"S",
							"S",
							"a",
						},
						{
							"S",
							"b",
							"S",
							"S",
						},
						{
							"a",
						},
					},
				},
				idGetter: newIDGetter(),
			},
			want: []models.Rule{
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"S",
							"[new_NT_0]",
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							"S",
							"a",
						},
					},
				},
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"S",
							"[new_NT_1]",
						},
					},
				},
				{
					NonTerminal: "[new_NT_1]",
					Rights: []models.ProductionBody{
						{
							"b",
							"[new_NT_2]",
						},
					},
				},
				{
					NonTerminal: "[new_NT_2]",
					Rights: []models.ProductionBody{
						{
							"S",
							"S",
						},
					},
				},
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							"a",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deleteLongPart(tt.args.rule, tt.args.idGetter)

			require.Equal(t, tt.want, result)
		})
	}
}

func Test_deleteChainRules(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "from itmo",
			input: "S -> B | dd\nB -> cc",
			want:  "S -> cc | dd\nB -> cc",
		},
		{
			name:  "my 1",
			input: "S -> B | E | D\nB -> cc\nE -> e\nD -> ddd",
			want:  "S -> cc | e | ddd\nB -> cc\nE -> e\nD -> ddd",
		},
		{
			name:  "itmo 2",
			input: "S -> B | a\nB -> C | b\nC -> dd | c",
			want:  "S -> dd | c | b | a\nB -> dd | c | b\nC -> dd | c",
		},
		{
			name:  "my 2",
			input: "S -> B\nB -> D\nD -> ddd",
			want:  "S -> ddd\nB -> ddd\nD -> ddd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := parser.New().Parse(tt.want, "S")
			input := parser.New().Parse(tt.input, "S")

			result := deleteChainRules(input)

			require.Equal(t, expected, result)
		})
	}
}

func Test_deleteNonGenerative(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "from itmo",
			input: "S -> Ac\nA->SD\nD->aD\nA->a",
			want:  "S -> Ac\nA -> a",
		},
		{
			name:  "from wiki",
			input: "S -> Bb | Ee\nE -> Ee\nB -> b",
			want:  "S -> Bb\nB -> b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := parser.New().Parse(tt.want, "S")
			input := parser.New().Parse(tt.input, "S")

			result := deleteNonGenerative(input)

			require.Equal(t, expected, result)
		})
	}
}

func Test_deleteNonReachable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "itmo",
			input: "S -> AB | CD\nA -> EF\nG -> AD\nC -> c",
			want:  "S -> AB | CD\nA -> EF\nC -> c",
		},
		{
			name:  "itmo",
			input: "S -> B\nC -> c\nD -> d\nE->e",
			want:  "S -> B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := parser.New().Parse(tt.want, "S")
			input := parser.New().Parse(tt.input, "S")

			result := deleteNonReachable(input)

			require.Equal(t, expected, result)
		})
	}
}

func Test_deletePairedTerminals(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *grammar.Grammar
	}{
		{
			name:  "my",
			input: "S -> aa",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_0]",
								"[NT_PT_0]",
							},
						},
					},
					"[NT_PT_0]": {
						NonTerminal: "[NT_PT_0]",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
				},
			},
		},
		{
			name:  "my",
			input: "S -> aa | bb | dC\nC -> cc",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_0]",
								"[NT_PT_0]",
							},
							{
								"[NT_PT_1]",
								"[NT_PT_1]",
							},
							{
								"[NT_PT_2]",
								"C",
							},
						},
					},
					"[NT_PT_0]": {
						NonTerminal: "[NT_PT_0]",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
					"[NT_PT_1]": {
						NonTerminal: "[NT_PT_1]",
						Rights: []models.ProductionBody{
							{
								"b",
							},
						},
					},
					"[NT_PT_2]": {
						NonTerminal: "[NT_PT_2]",
						Rights: []models.ProductionBody{
							{
								"d",
							},
						},
					},
					"C": {
						NonTerminal: "C",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_3]",
								"[NT_PT_3]",
							},
						},
					},
					"[NT_PT_3]": {
						NonTerminal: "[NT_PT_3]",
						Rights: []models.ProductionBody{
							{
								"c",
							},
						},
					},
				},
			},
		},
		{
			name:  "my",
			input: "S -> a",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
				},
			},
		},
		{
			name:  "my",
			input: "S -> ab",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_0]",
								"[NT_PT_1]",
							},
						},
					},
					"[NT_PT_0]": {
						NonTerminal: "[NT_PT_0]",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
					"[NT_PT_1]": {
						NonTerminal: "[NT_PT_1]",
						Rights: []models.ProductionBody{
							{
								"b",
							},
						},
					},
				},
			},
		},
		{
			name:  "my",
			input: "S -> ab | cc | cD\nD -> aa",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_0]",
								"[NT_PT_1]",
							},
							{
								"[NT_PT_2]",
								"[NT_PT_2]",
							},
							{
								"[NT_PT_3]",
								"D",
							},
						},
					},
					"[NT_PT_0]": {
						NonTerminal: "[NT_PT_0]",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
					"[NT_PT_1]": {
						NonTerminal: "[NT_PT_1]",
						Rights: []models.ProductionBody{
							{
								"b",
							},
						},
					},
					"[NT_PT_2]": {
						NonTerminal: "[NT_PT_2]",
						Rights: []models.ProductionBody{
							{
								"c",
							},
						},
					},
					"[NT_PT_3]": {
						NonTerminal: "[NT_PT_3]",
						Rights: []models.ProductionBody{
							{
								"c",
							},
						},
					},
					"D": {
						NonTerminal: "D",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_4]",
								"[NT_PT_4]",
							},
						},
					},
					"[NT_PT_4]": {
						NonTerminal: "[NT_PT_4]",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parser.New().Parse(tt.input, "S")

			result := deletePairedTerminals(input)

			require.Equal(t, tt.want, result)
		})
	}
}

func TestCNF_ToCNF(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *grammar.Grammar
	}{
		{
			name:  "1",
			input: "S -> cA | dA | cB | eB\nA -> a\nB -> b",
			want: &grammar.Grammar{
				Start: "S",
				Grammar: map[string]models.Rule{
					"S": {
						NonTerminal: "S",
						Rights: []models.ProductionBody{
							{
								"[NT_PT_0]",
								"A",
							},
							{
								"[NT_PT_1]",
								"A",
							},
							{
								"[NT_PT_2]",
								"B",
							},
							{
								"[NT_PT_3]",
								"B",
							},
						},
					},
					"A": {
						NonTerminal: "A",
						Rights: []models.ProductionBody{
							{
								"a",
							},
						},
					},
					"B": {
						NonTerminal: "B",
						Rights: []models.ProductionBody{
							{
								"b",
							},
						},
					},
					"[NT_PT_0]": {
						NonTerminal: "[NT_PT_0]",
						Rights: []models.ProductionBody{
							{
								"c",
							},
						},
					},
					"[NT_PT_1]": {
						NonTerminal: "[NT_PT_1]",
						Rights: []models.ProductionBody{
							{
								"d",
							},
						},
					},
					"[NT_PT_2]": {
						NonTerminal: "[NT_PT_2]",
						Rights: []models.ProductionBody{
							{
								"c",
							},
						},
					},
					"[NT_PT_3]": {
						NonTerminal: "[NT_PT_3]",
						Rights: []models.ProductionBody{
							{
								"e",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := &CNF{}

			input := parser.New().Parse(tt.input, "S")

			result := cnf.ToCNF(input)

			require.Equal(t, tt.want, result)
		})
	}
}
