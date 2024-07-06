package models

import (
	"errors"
	"fmt"
)

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

type TypeBindings map[string]Type

type Expression interface {
	Evaluate(Bindings) (any, *InterpreterError)
	Type(TypeBindings) (Type, *InterpreterError)
	SourceLocation() SourceLocation
}

type Function interface {
	Call([]any) (any, error)
}

type Type interface {
	fmt.Stringer
	Name() (string, error)
	IsSuperTo(Type) bool
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

func (t PrimitiveType) String() string {
	name, err := t.Name()
	if err != nil {
		return "unknown"
	}
	return name
}

func (t PrimitiveType) IsSuperTo(other Type) bool {
	if t == other || t == PrimitiveTypeAny {
		return true
	}
	return false
}
