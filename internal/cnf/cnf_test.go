package cnf

import (
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
							Body: []string{
								"S",
								"S",
								"a",
							},
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
							Body: []string{
								"S",
								"[new_NT_0]",
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							Body: []string{
								"S",
								"a",
							},
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
							Body: []string{

								"S",
								"S",
								"a",
							},
						},
						{
							Body: []string{
								"S",
								"b",
								"S",
								"S",
							},
						},
						{
							Body: []string{
								"a",
							},
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
							Body: []string{
								"S",
								"[new_NT_0]",
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							Body: []string{
								"S",
								"a",
							},
						},
					},
				},
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							Body: []string{
								"S",
								"[new_NT_1]",
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_1]",
					Rights: []models.ProductionBody{
						{
							Body: []string{
								"b",
								"[new_NT_2]",
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_2]",
					Rights: []models.ProductionBody{
						{
							Body: []string{
								"S",
								"S",
							},
						},
					},
				},
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							Body: []string{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := parser.New().Parse(tt.want)
			input := parser.New().Parse(tt.input)

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
			expected := parser.New().Parse(tt.want)
			input := parser.New().Parse(tt.input)

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
			expected := parser.New().Parse(tt.want)
			input := parser.New().Parse(tt.input)

			result := deleteNonReachable(input)

			require.Equal(t, expected, result)
		})
	}
}
