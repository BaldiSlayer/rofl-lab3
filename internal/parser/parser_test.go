package parser

import (
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
		want SymbolsBtw
	}{
		{
			name: "first",
			args: args{
				input: "aba]",
			},
			want: SymbolsBtw{
				s: "[aba]",
			},
		},
		{
			name: "second",
			args: args{
				input: "fasdfasdfsda]",
			},
			want: SymbolsBtw{
				s: "[fasdfasdfsda]",
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
		want []ProductionBody
	}{
		{

			name: "simple ",
			args: args{
				s: "SSa|SbSS|a",
			},
			want: []ProductionBody{
				{
					body: []SymbolsBtw{
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
					body: []SymbolsBtw{
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
					body: []SymbolsBtw{
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
			want: []ProductionBody{
				{
					body: []SymbolsBtw{
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
		want Rule
	}{
		{
			name: "first",
			args: args{
				s: "S  -> SSa  |SbSS| a",
			},
			want: Rule{
				nonTerminal: "S",
				rights: []ProductionBody{
					{
						body: []SymbolsBtw{
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
						body: []SymbolsBtw{
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
						body: []SymbolsBtw{
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
			want: Rule{
				nonTerminal: "S",
				rights: []ProductionBody{
					{
						body: []SymbolsBtw{
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
