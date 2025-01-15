package main

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cnf"
	"github.com/BaldiSlayer/rofl-lab3/internal/fuzzer"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
)

func main() {
	p := parser.New()

	fuzz := fuzzer.New(
		p,
		&cnf.CNF{},
		&bigramms.Bigramms{},
	)

	results := fuzz.Generate(100, 0.1, "S")
	for _, line := range results {
		fmt.Println(line)
	}
}
