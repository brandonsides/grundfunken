package models

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
	Type() (any, *InterpreterError)
	SourceLocation() SourceLocation
}

type Function interface {
	Call([]any) (any, error)
}

func TypeInt(v any) bool {
	_, ok := v.(int)
	return ok
}

func TypeString(v any) bool {
	_, ok := v.(string)
	return ok
}

func TypeBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

func TypeObject(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

func TypeFunction(v any) bool {
	_, ok := v.(Function)
	return ok
}

func TypeArray(v any) bool {
	_, ok := v.([]any)
	return ok
}
