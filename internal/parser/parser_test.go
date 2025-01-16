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
		want string
	}{
		{
			name: "first",
			args: args{
				input: "[aba]",
			},
			want: "[aba]",
		},
		{
			name: "second",
			args: args{
				input: "[fasdfasdfsda]",
			},
			want: "[fasdfasdfsda]",
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
		want []grammar.ProductionBody
	}{
		{

			name: "simple ",
			args: args{
				s: "SSa|SbSS|a",
			},
			want: []grammar.ProductionBody{
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
		{

			name: "many S",
			args: args{
				s: "SSSSSSSS",
			},
			want: []grammar.ProductionBody{
				{
					"S",
					"S",
					"S",
					"S",
					"S",
					"S",
					"S",
					"S",
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
		want grammar.Rule
	}{
		{
			name: "first",
			args: args{
				s: "S  -> SSa  |SbSS| a",
			},
			want: grammar.Rule{
				NonTerminal: "S",
				Rights: []grammar.ProductionBody{
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
		},
		{
			name: "first",
			args: args{
				s: "S  ->    SSSSSSSS",
			},
			want: grammar.Rule{
				NonTerminal: "S",
				Rights: []grammar.ProductionBody{
					{
						"S",
						"S",
						"S",
						"S",
						"S",
						"S",
						"S",
						"S",
					},
				},
			},
		},
		{
			name: "name",
			args: args{
				s: "S -> [order66]ab",
			},
			want: grammar.Rule{
				NonTerminal: "S",
				Rights: []grammar.ProductionBody{
					{
						"[order66]",
						"a",
						"b",
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
