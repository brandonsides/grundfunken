package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
)

type FieldAccessExpression struct {
	Object   models.Expression
	Field    string
	fieldLoc models.SourceLocation
}

func (fae *FieldAccessExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	t, err := fae.Object.Type(tb)
	if err != nil {
		return nil, err
	}

	tObj, ok := t.(types.ObjectType)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("cannot access field on type %s", t.String()),
			SourceLocation: fae.Object.SourceLocation(),
		}
	}

	fieldType, ok := tObj.Fields[fae.Field]
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("field %s not found on type %s", fae.Field, t.String()),
			SourceLocation: fae.fieldLoc,
		}
	}

	return fieldType, nil
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
