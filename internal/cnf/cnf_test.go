package cnf

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/grammar"
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
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
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"S",
								},
								{
									"a",
								},
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
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"[new_NT_0]",
								},
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"a",
								},
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
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"S",
								},
								{
									"a",
								},
							},
						},
						{
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"b",
								},
								{
									"S",
								},
								{
									"S",
								},
							},
						},
						{
							Body: []models.SymbolsBtw{
								{
									"a",
								},
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
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"[new_NT_0]",
								},
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_0]",
					Rights: []models.ProductionBody{
						{
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"a",
								},
							},
						},
					},
				},
				{
					NonTerminal: "B",
					Rights: []models.ProductionBody{
						{
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"[new_NT_1]",
								},
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_1]",
					Rights: []models.ProductionBody{
						{
							Body: []models.SymbolsBtw{
								{
									"b",
								},
								{
									"[new_NT_2]",
								},
							},
						},
					},
				},
				{
					NonTerminal: "[new_NT_2]",
					Rights: []models.ProductionBody{
						{
							Body: []models.SymbolsBtw{
								{
									"S",
								},
								{
									"S",
								},
							},
						},
					},
				},
				{
					NonTerminal: "B",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deleteLongPart(tt.args.rule, tt.args.idGetter)

			require.Equal(t, tt.want, result)
		})
	}
}

func Test_deleteChainRules_1(t *testing.T) {
	expected := &grammar.Grammar{
		Start: "A",
		Grammar: map[string]models.Rule{
			"A": {
				NonTerminal: "A",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"c",
							},
							{
								"c",
							},
						},
					},
					{
						Body: []models.SymbolsBtw{
							{
								"d",
							},
							{
								"d",
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
								"c",
							},
							{
								"c",
							},
						},
					},
				},
			},
		},
	}

	args := &grammar.Grammar{
		Start: "A",
		Grammar: map[string]models.Rule{
			"A": {
				NonTerminal: "A",
				Rights: []models.ProductionBody{
					{
						Body: []models.SymbolsBtw{
							{
								"B",
							},
						},
					},
					{
						Body: []models.SymbolsBtw{
							{
								"d",
							},
							{
								"d",
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
								"c",
							},
							{
								"c",
							},
						},
					},
				},
			},
		},
	}

	result := deleteChainRules(args)

	require.Equal(t, expected, result, "")
}
