package models

import "errors"

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
	Type() (Type, *InterpreterError)
	SourceLocation() SourceLocation
}

type Function interface {
	Call([]any) (any, error)
}

type Type interface {
	Name() (string, error)
}

type PrimitiveType uint8

const (
	PrimitiveTypeInt PrimitiveType = iota
	PrimitiveTypeString
	PrimitiveTypeBool
	PrimitiveTypeList
	PrimitiveTypeObject
	PrimitiveTypeFunction
	PrimitiveTypeAny
)

func (t PrimitiveType) Name() (string, error) {
	switch t {
	case PrimitiveTypeInt:
		return "int", nil
	case PrimitiveTypeString:
		return "string", nil
	case PrimitiveTypeBool:
		return "bool", nil
	case PrimitiveTypeList:
		return "list", nil
	case PrimitiveTypeObject:
		return "object", nil
	case PrimitiveTypeFunction:
		return "function", nil
	default:
		return "", errors.New("unknown primitive type")
	}
}
