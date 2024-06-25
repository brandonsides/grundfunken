package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/parser"
	"github.com/brandonksides/grundfunken/tokens"
)

func main() {
	var inputFilePath string
	flag.StringVar(&inputFilePath, "input", "", "Path to the input file")
	flag.Parse()

	var input io.ReadCloser
	if inputFilePath == "" {
		input = os.Stdin
	} else {
		var err error
		input, err = os.Open(inputFilePath)
		if err != nil {
			panic(fmt.Errorf("failed to open the file at the provided path: %w", err))
		}
	}
	defer input.Close()

	// hold all the input lines in memory
	// so we can report errors with context
	var lines []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading input: %v\n", err)
		return
	}

	result, err := interpret(lines)
	if err != nil {
		report(err, lines)
		return
	}

	fmt.Printf("Result: %v\n", result)
}

func interpret(lines []string) (any, *models.InterpreterError) {
	// split input into "tokens", which are the smallest
	// meaningful units of the language: words, numbers,
	// punctuation, etc.
	toks, err := tokens.Tokenize(lines)
	if err != nil {
		return nil, err
	}

	// parse the tokens into an "expression", which is a
	// tree-like structure that represents the semantic
	// relationships between the tokens
	expression, rest, err := parser.ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	if len(rest) != 0 {
		return nil, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: rest[0].SourceLocation,
		}
	}

	// evaluate the expression to get the final result
	// with the top-level bindings for certain builtin
	// identifiers
	return expression.Evaluate(builtins)
}

func report(err error, lines []string) {
	fmt.Print("Error: ")
	reportHelper(err, lines)
}

func reportHelper(err error, lines []string) {
	interpreterErr, ok := err.(*models.InterpreterError)
	if !ok {
		fmt.Println(err.Error())
		fmt.Println()
		return
	}

	if interpreterErr.Underlying != nil {
		reportHelper(interpreterErr.Underlying, lines)
	}

	fmt.Println(interpreterErr.Error())
	fmt.Println(lines[interpreterErr.SourceLocation.LineNumber])
	fmt.Println(underlineError(interpreterErr.SourceLocation.ColumnNumber))
}

func underlineError(col int) string {
	underline := ""
	for i := 0; i < col; i++ {
		underline += " "
	}
	underline += "^-here\n"
	return underline
}
