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

type ExpFunction struct {
	Args     []string
	Bindings Bindings
	Exp      Expression
}

func (f ExpFunction) Call(args []any) (any, *InterpreterError) {
	if len(args) != len(f.Args) {
		return nil, &InterpreterError{
			Err: fmt.Errorf("expected %d arguments, got %d", len(f.Args), len(args)),
		}
	}
	newBindings := make(Bindings)
	for k, v := range f.Bindings {
		newBindings[k] = v
	}
	for i, arg := range f.Args {
		newBindings[arg] = args[i]
	}
	return f.Exp.Evaluate(newBindings)
}

type Function interface {
	Call([]any) (any, *InterpreterError)
}
