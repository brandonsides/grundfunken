package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
)

type FieldAccessExpression struct {
	Object models.Expression
	Field  string
	loc    models.SourceLocation
}

func (fae *FieldAccessExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	obj, err := fae.Object.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	objMap, ok := obj.(map[string]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("expected object"),
			SourceLocation: fae.Object.SourceLocation(),
		}
	}

	if val, ok := objMap[fae.Field]; ok {
		return val, nil
	}

	return nil, &models.InterpreterError{
		Err:            errors.New("field not found"),
		SourceLocation: fae.loc,
	}
}

func (fae *FieldAccessExpression) SourceLocation() models.SourceLocation {
	return fae.loc
}