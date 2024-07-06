package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/parser"
	"github.com/brandonksides/grundfunken/tokens"
)

func main() {
	var inputFilePath string
	flag.StringVar(&inputFilePath, "input", "", "Path to the input file")
	flag.Parse()

	result, lines, err := interpret(inputFilePath)
	if err != nil {
		report(err, lines)
		return
	}

	fmt.Printf("Result: %v\n", result)
}

func interpret(inputFilePath string) (any, map[string][]string, error) {
	var input io.ReadCloser

	var fileName string
	if inputFilePath == "" {
		input = os.Stdin
		fileName = "stdin"
	} else {
		var err error
		input, err = os.Open(inputFilePath)
		if err != nil {
			panic(fmt.Errorf("failed to open the file at the provided path: %w", err))
		}

		splitPath := strings.Split(inputFilePath, "/")
		oldDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("failed to get the current working directory: %w", err))
		}
		os.Chdir(strings.Join(splitPath[:len(splitPath)-1], "/"))
		defer os.Chdir(oldDir)
		fileName = splitPath[len(splitPath)-1]
	}
	defer input.Close()

	// hold all the input mainLines in memory
	// so we can report errors with context
	lines := map[string][]string{fileName: make([]string, 0)}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines[fileName] = append(lines[fileName], scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading input: %v\n", err)
		return nil, lines, err
	}
	// split input into "tokens", which are the smallest
	// meaningful units of the language: words, numbers,
	// punctuation, etc.
	toks, err := tokens.Tokenize(fileName, lines[fileName])
	if err != nil {
		return nil, lines, err
	}

	// parse the tokens into an "expression", which is a
	// tree-like structure that represents the semantic
	// relationships between the tokens
	expression, err := parser.ParseExpression(toks)
	if err != nil {
		return nil, lines, err
	}

	tok, ok := toks.Peek()
	if ok {
		return nil, lines, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: tok.SourceLocation,
		}
	}

	// evaluate the expression to get the final result
	// with the top-level bindings for certain builtin
	// identifiers
	builtins["import"] = &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, error) {
			path := args[0].(string)

			ret, newLines, err := interpret(path)
			for fileName, fileLines := range newLines {
				lines[fileName] = fileLines
			}
			return ret, err
		},
	}

	var builtinTypes = map[string]models.Type{}
	for k := range builtins {
		builtinTypes[k] = models.PrimitiveTypeFunction
	}

	_, err = expression.Type(builtinTypes)
	if err != nil {
		return nil, lines, err
	}

	ret, err := expression.Evaluate(builtins)
	if err != nil {
		return ret, lines, err
	}
	return ret, lines, nil
}

func report(err error, lines map[string][]string) {
	fmt.Print("Error: ")
	reportHelper(err, lines)
}

func reportHelper(err error, lines map[string][]string) {
	interpreterErr, ok := err.(*models.InterpreterError)
	if !ok {
		fmt.Println(err.Error())
		fmt.Println()
		return
	}

	if interpreterErr.Underlying != nil {
		reportHelper(interpreterErr.Underlying, lines)
	}

	fmt.Printf("in file %s at line %d, column %d: %s\n", interpreterErr.SourceLocation.File, interpreterErr.SourceLocation.LineNumber+1, interpreterErr.SourceLocation.ColumnNumber+1, interpreterErr.Error())
	fmt.Println(lines[interpreterErr.SourceLocation.File][interpreterErr.SourceLocation.LineNumber])
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
