package parser

import "github.com/brandonksides/grundfunken/models"

type LiteralExpression struct {
	val any
	loc models.SourceLocation
}

func (le *LiteralExpression) Type(tb models.TypeBindings) (models.Type, *models.InterpreterError) {
	switch le.val.(type) {
	case bool:
		return models.PrimitiveTypeBool, nil
	case int:
		return models.PrimitiveTypeInt, nil
	case string:
		return models.PrimitiveTypeString, nil
	default:
		return models.PrimitiveTypeAny, nil
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
