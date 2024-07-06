package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
)

type FieldAccessExpression struct {
	Object   models.Expression
	Field    string
	fieldLoc models.SourceLocation
}

func (fae *FieldAccessExpression) Type(tb models.TypeBindings) (models.Type, *models.InterpreterError) {
	t, err := fae.Object.Type(tb)
	if err != nil {
		return nil, err
	}

	if t != models.PrimitiveTypeObject {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("cannot access field on type %s", t.String()),
			SourceLocation: fae.Object.SourceLocation(),
		}
	}

	return models.PrimitiveTypeAny, nil
}

func (fae *FieldAccessExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	obj, err := fae.Object.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	objMap, ok := obj.(map[string]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected object; got %v", obj),
			SourceLocation: fae.Object.SourceLocation(),
		}
	}

	if val, ok := objMap[fae.Field]; ok {
		return val, nil
	}

	return nil, &models.InterpreterError{
		Message:        "field not found",
		SourceLocation: fae.fieldLoc,
	}
}

func (fae *FieldAccessExpression) SourceLocation() models.SourceLocation {
	return fae.Object.SourceLocation()
}
