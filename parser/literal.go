package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
)

type LiteralExpression struct {
	val any
	loc models.SourceLocation
}

func (le *LiteralExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	switch le.val.(type) {
	case bool:
		return types.PrimitiveTypeBool, nil
	case int:
		return types.PrimitiveTypeInt, nil
	case string:
		return types.PrimitiveTypeString, nil
	default:
		return types.PrimitiveTypeAny, nil
	}
}

func (le *LiteralExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	if le == nil {
		return nil, nil
	}
	return le.val, nil
}

func (le *LiteralExpression) SourceLocation() models.SourceLocation {
	return le.loc
}
