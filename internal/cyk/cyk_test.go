package cyk

import (
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCYK_Check_1(t *testing.T) {
	input := parser.New().Parse("S -> a")

	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "1",
			args: "a",
			want: true,
		},
		{
			name: "2",
			args: "aa",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := New(input).Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}

func TestCYK_Check_2(t *testing.T) {
	input := parser.New().Parse("S -> AA\nA -> a")

	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "1",
			args: "a",
			want: false,
		},
		{
			name: "2",
			args: "aa",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := New(input).Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}

func TestCYK_Check_PSP(t *testing.T) {
	input := parser.New().Parse("S -> BB | CD\nB -> BB | CD\nC -> a\nD -> BE | b\nE -> b")

	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "1",
			args: "dd",
			want: false,
		},
		{
			name: "2",
			args: "dd",
			want: false,
		},
		{
			name: "3",
			args: "abb",
			want: false,
		},
		{
			name: "3",
			args: "ab",
			want: true,
		},
		{
			name: "4",
			args: "aabb",
			want: true,
		},
		{
			name: "5",
			args: "abababababababababababababaabababababababab",
			want: false,
		},
		{
			name: "6",
			args: "abaabb",
			want: true,
		},
		{
			name: "7",
			args: "aabbab",
			want: true,
		},
		{
			name: "8",
			args: "abaabbab",
			want: true,
		},
		{
			name: "9",
			args: "ababb",
			want: false,
		},
		{
			name: "10",
			args: "abbb",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := New(input).Check(tt.args)

			require.Equal(t, tt.want, res)
		})
	}
}
