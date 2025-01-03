package parser

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseBetweenBrackets(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want models.SymbolsBtw
	}{
		{
			name: "first",
			args: args{
				input: "aba]",
			},
			want: models.SymbolsBtw{
				S: "[aba]",
			},
		},
		{
			name: "second",
			args: args{
				input: "fasdfasdfsda]",
			},
			want: models.SymbolsBtw{
				S: "[fasdfasdfsda]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Parser

			result := p.parseBetweenBrackets(tt.args.input)

			require.Equal(t, tt.want, result)
		})
	}
}

func Test_parseRight(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []models.ProductionBody
	}{
		{

			name: "simple ",
			args: args{
				s: "SSa|SbSS|a",
			},
			want: []models.ProductionBody{
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
		{

			name: "many S",
			args: args{
				s: "SSSSSSSS",
			},
			want: []models.ProductionBody{
				{
					Body: []models.SymbolsBtw{
						{
							"S",
						},
						{
							"S",
						},
						{
							"S",
						},
						{
							"S",
						},
						{
							"S",
						},
						{
							"S",
						},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Parser

			got := p.parseRight(tt.args.s)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestParser_parseLine(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want models.Rule
	}{
		{
			name: "first",
			args: args{
				s: "S  -> SSa  |SbSS| a",
			},
			want: models.Rule{
				NonTerminal: "S",
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
		},
		{
			name: "first",
			args: args{
				s: "S  ->    SSSSSSSS",
			},
			want: models.Rule{
				NonTerminal: "S",
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
								"S",
							},
							{
								"S",
							},
							{
								"S",
							},
							{
								"S",
							},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}

			got := p.parseLine(tt.args.s)

			require.Equal(t, tt.want, got)
		})
	}
}
