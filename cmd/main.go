package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BaldiSlayer/rofl-lab3/internal/bigramms"
	"github.com/BaldiSlayer/rofl-lab3/internal/cnf"
	"github.com/BaldiSlayer/rofl-lab3/internal/fuzzer"
	"github.com/BaldiSlayer/rofl-lab3/internal/parser"
)

var (
	testsCount         = flag.Int("tests_count", 10, "Number of tests")
	startSymbol        = flag.String("start_symbol", "S", "Starting symbol")
	breakProb          = flag.Float64("break_prob", 0.1, "Probability of break (0.0 - 1.0)")
	terminalAddingProb = flag.Float64("terminal_adding_prob", 0.1, "Probability of terminal addition (0.0 - 1.0)")
	help               = flag.Bool("help", false, "Show help message")
)

func printHelp() {
	fmt.Println("Usage: go run main.go [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  -tests_count int        Number of tests")
	fmt.Println("  -start_symbol string    Starting symbol")
	fmt.Println("  -break_prob float        Probability of break (0.0 - 1.0)")
	fmt.Println("  -terminal_adding_prob float Probability of terminal addition (0.0 - 1.0)")
}

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
	flag.Parse()

	if *help {
		printHelp()

		os.Exit(0)
	}

	if *testsCount <= 0 || *startSymbol == "" || *breakProb < 0.0 ||
		*breakProb > 1.0 || *terminalAddingProb < 0.0 || *terminalAddingProb > 1.0 {
		fmt.Println("Invalid arguments.")
		printHelp()

		os.Exit(1)
	}

	input := inputLines()

	fuzz := fuzzer.New(
		input,
		parser.New(),
		&cnf.CNF{},
		&bigramms.Bigramms{},
		*startSymbol,
	)

	results := fuzz.Generate(*testsCount, *breakProb, *terminalAddingProb)
	for _, line := range results {
		fmt.Println(line)
	}
}
