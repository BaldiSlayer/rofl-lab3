package main

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cnf"
	"github.com/BaldiSlayer/rofl-lab3/internal/fuzzer"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
)

const (
	n = 100
	// todo rename
	someValue = 0.1

	startSmb = "S"
)

func main() {
	input := ""

	fuzz := fuzzer.New(
		input,
		parser.New(),
		&cnf.CNF{},
		&bigramms.Bigramms{},
	)

	results := fuzz.Generate(n, someValue, startSmb)
	for _, line := range results {
		fmt.Println(line)
	}
}
