package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cnf"
	"github.com/BaldiSlayer/rofl-lab3/internal/fuzzer"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
)

const (
	n = 10
	// todo rename
	someValue = 0.1

	startSmb = "S"
)

func inputLines() string {
	lines := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка ввода:", err)

		return ""
	}

	var sb strings.Builder

	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}

func main() {
	input := inputLines()

	fuzz := fuzzer.New(
		input,
		parser.New(),
		&cnf.CNF{},
		&bigramms.Bigramms{},
		"S",
	)

	results := fuzz.Generate(n, someValue, startSmb)
	for _, line := range results {
		fmt.Println(line)
	}
}
