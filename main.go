package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/brandonksides/crazy-types/models"
	"github.com/brandonksides/crazy-types/parser"
	"github.com/brandonksides/crazy-types/tokens"
)

type logger struct{}

func (l logger) Errorf(err models.InterpreterError) {
	fmt.Println(err.Error())
}

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
		report(*err, lines)
		return
	}

	fmt.Printf("Result: %v", result)
}

func interpret(lines []string) (any, *models.InterpreterError) {
	toks, err := tokens.Tokenize(lines)
	if err != nil {
		return nil, err
	}

	expression, rest, err := parser.ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	if len(rest) != 0 {
		return nil, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	return expression.Evaluate(models.Bindings{})
}

func report(err models.InterpreterError, lines []string) {
	fmt.Println(err.Error())
	fmt.Println(lines[err.SourceLocation.LineNumber])
	fmt.Println(underlineError(err.SourceLocation.ColumnNumber))
}

func underlineError(col int) string {
	underline := ""
	for i := 0; i < col; i++ {
		underline += " "
	}
	underline += "^-here\n"
	return underline
}
