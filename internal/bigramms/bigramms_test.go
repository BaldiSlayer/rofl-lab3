package bigramms

import (
	"testing"

	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	dest := map[string]struct{}{"a": {}, "b": {}}
	src := map[string]struct{}{"b": {}, "c": {}}
	expected := map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithEmptySrc(t *testing.T) {
	dest := map[string]struct{}{"x": {}, "y": {}}
	src := map[string]struct{}{}
	expected := map[string]struct{}{
		"x": {},
		"y": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithEmptyDest(t *testing.T) {
	dest := map[string]struct{}{}
	src := map[string]struct{}{"m": {}, "n": {}}
	expected := map[string]struct{}{
		"m": {},
		"n": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func TestUnionWithIdenticalMaps(t *testing.T) {
	dest := map[string]struct{}{"p": {}}
	src := map[string]struct{}{"p": {}}
	expected := map[string]struct{}{
		"p": {},
	}

	result := union(dest, src)

	assert.Equal(t, expected, result)
}

func Test_makeFirst(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]map[string]struct{}
	}{
		{
			name:  "1",
			input: "S -> AB\nA -> a\nB -> b",
			want: map[string]map[string]struct{}{
				"A": {
					"a": {},
				},
				"B": {
					"b": {},
				},
				"S": {
					"a": {},
				},
			},
		},
		{
			name:  "2",
			input: "S -> a | b\nA -> a\nB -> b",
			want: map[string]map[string]struct{}{
				"S": {
					"a": {},
					"b": {},
				},
				"A": {
					"a": {},
				},
				"B": {
					"b": {},
				},
			},
		},
		{
			name:  "3",
			input: "S -> B | A\nA -> a\nB -> b",
			want: map[string]map[string]struct{}{
				"S": {
					"a": {},
					"b": {},
				},
				"A": {
					"a": {},
				},
				"B": {
					"b": {},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parser.New().Parse(tt.input, "S")

			result := makeFirst(input)

			require.Equal(t, tt.want, result)
		})
	}
}

func Test_makeLast(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]map[string]struct{}
	}{
		{
			name:  "1",
			input: "S -> AB\nA -> a\nB -> b",
			want: map[string]map[string]struct{}{
				"A": {
					"a": {},
				},
				"B": {
					"b": {},
				},
				"S": {
					"b": {},
				},
			},
		},
		{
			name:  "2",
			input: "S -> a | b\nA -> a\nB -> b",
			want: map[string]map[string]struct{}{
				"S": {
					"a": {},
					"b": {},
				},
				"A": {
					"a": {},
				},
				"B": {
					"b": {},
				},
			},
		},
		{
			name:  "3",
			input: "S -> B | A\nA -> ab\nB -> bc",
			want: map[string]map[string]struct{}{
				"S": {
					"c": {},
					"b": {},
				},
				"A": {
					"b": {},
				},
				"B": {
					"c": {},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parser.New().Parse(tt.input, "S")

			result := makeLast(input)

			require.Equal(t, tt.want, result)
		})
	}
}
