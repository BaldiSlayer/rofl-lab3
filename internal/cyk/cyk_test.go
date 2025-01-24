package cyk

import (
	"testing"

	"github.com/BaldiSlayer/rofl-lab3/internal/parser"

	"github.com/stretchr/testify/require"
)

func TestCYK_Check_1(t *testing.T) {
	input := parser.New().Parse("S -> a", "S")

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
	input := parser.New().Parse("S -> AA\nA -> a", "S")

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
	input := parser.New().Parse("S -> BB | CD\nB -> BB | CD\nC -> a\nD -> BE | b\nE -> b", "S")

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

func TestCYK_Check_3(t *testing.T) {
	input := parser.New().Parse(
		"[new_NT_0] -> [NT_PT_0][NT_PT_1]\n[order66] -> a\n[NT_PT_0] -> a\n"+
			"S -> [order66][new_NT_0] | a\n[NT_PT_1] -> b",
		"S",
	)

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
			name: "1",
			args: "b",
			want: false,
		},
		{
			name: "1",
			args: "d",
			want: false,
		},
		{
			name: "3",
			args: "aab",
			want: true,
		},
		{
			name: "4",
			args: "abababababbabbababbabbabbabbbabbabbabab",
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

func TestCYK_Check_4(t *testing.T) {
	input := parser.New().Parse(
		"S -> [NT_PT_a][new_NT_0] | [NT_ALONE_c] | [NT_PT_a][new_NT_1] | [NT_ALONE_d]\n"+
			"[new_NT_0] -> C[NT_PT_b]\n"+
			"C -> [NT_PT_a][new_NT_0] | [NT_ALONE_c]\n"+
			"D -> [NT_PT_a][new_NT_1] | [NT_ALONE_d]\n"+
			"[NT_PT_a] -> [NT_ALONE_a]\n"+
			"[NT_ALONE_a] -> a\n"+
			"[NT_ALONE_d] -> d\n"+
			"[NT_ALONE_b] -> b\n"+
			"[new_NT_2] -> [NT_PT_b][NT_PT_b]\n"+
			"[new_NT_1] -> D[new_NT_2]\n"+
			"[NT_PT_b] -> [NT_ALONE_b]\n[NT_ALONE_c] -> c\n",
		"S",
	)

	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "1",
			args: "c",
			want: true,
		},
		{
			name: "2",
			args: "d",
			want: true,
		},
		{
			name: "3",
			args: "ddd",
			want: false,
		},
		{
			name: "4",
			args: "a",
			want: false,
		},
		{
			name: "4",
			args: "adbb",
			want: true,
		},
		{
			name: "5",
			args: "aaaaacbbbbb",
			want: true,
		},
		{
			name: "5",
			args: "aaaaadbbbbbbbbbb",
			want: true,
		},
		{
			name: "5",
			args: "e",
			want: false,
		},
		{
			name: "9",
			args: "aaaadbbbbbbbb",
			want: true,
		},
		{
			name: "9",
			args: "aaaadbbbbbbbcb",
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
