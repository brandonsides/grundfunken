package models

import "github.com/brandonksides/grundfunken/models/types"

type InterpreterError struct {
	Message        string
	Underlying     error
	SourceLocation SourceLocation
}

type SourceLocation struct {
	File         string
	LineNumber   int
	ColumnNumber int
}

func (e *InterpreterError) Error() string {
	return e.Message
}

type Bindings map[string]any

type Expression interface {
	Evaluate(Bindings) (any, *InterpreterError)
	Type(types.TypeBindings) (types.Type, *InterpreterError)
	SourceLocation() SourceLocation
}
