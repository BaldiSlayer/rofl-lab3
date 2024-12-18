package grammar

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
			in := make(chan byte, 0)
			out := make(chan SymbolsBtw, 1)

			go func() {
				for _, nt := range parseBetweenBrackets(in, '[') {
					out <- nt
				}
			}()

			for _, hui := range tt.args.input {
				in <- byte(hui)
			}

			close(in)

			result := <-out

			require.Equal(t, tt.want, result)
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want []SymbolsBtw
	}{
		{
			name: "normal",
			args: args{
				input: "a[aba]b",
			},
			want: []SymbolsBtw{
				{
					s: "a",
				},
				{
					s: "[aba]",
				},
				{
					s: "b",
				},
			},
		},
		{
			name: "T in brackets (it is NT)",
			args: args{
				input: "[a][aba]b",
			},
			want: []SymbolsBtw{
				{
					s: "[a]",
				},
				{
					s: "[aba]",
				},
				{
					s: "b",
				},
			},
		},
		{
			name: "three T and one NT",
			args: args{
				input: "acd[BC9]",
			},
			want: []SymbolsBtw{
				{
					s: "a",
				},
				{
					s: "c",
				},
				{
					s: "d",
				},
				{
					s: "[BC9]",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan byte, 0)
			out := make(chan SymbolsBtw)

			go parse(in, out)

			go func() {
				for _, hui := range tt.args.input {
					in <- byte(hui)
				}

				close(in)
			}()

			body := make([]SymbolsBtw, 0)

			for c := range out {
				body = append(body, c)
			}

			require.Equal(t, tt.want, body)
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
			got := parseRight(tt.args.s)

			require.Equal(t, tt.want, got)
		})
	}
}
