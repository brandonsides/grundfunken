package parser

import (
	"errors"

	"github.com/brandonksides/grundfunken/models"
)

type IdentifierExpression struct {
	name string
	loc  models.SourceLocation
}

func (ie *IdentifierExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret, ok := map[string]any(bindings)[ie.name]
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("cannot evaluate unbound identifier"),
			SourceLocation: ie.loc,
		}
	}

	return ret, nil
}

func (ie *IdentifierExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}
