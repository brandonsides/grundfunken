package models

import "fmt"

type InterpreterError struct {
	Err            error
	SourceLocation SourceLocation
}

type SourceLocation struct {
	LineNumber   int
	ColumnNumber int
}

func (e *InterpreterError) Error() string {
	return fmt.Sprintf("error at line %d, column %d: %v", e.SourceLocation.LineNumber+1, e.SourceLocation.ColumnNumber+1, e.Err)
}

type Bindings map[string]any

type Expression interface {
	Evaluate(Bindings) (any, *InterpreterError)
	SourceLocation() SourceLocation
}

type Function struct {
	Args     []string
	Bindings Bindings
	Exp      Expression
}
