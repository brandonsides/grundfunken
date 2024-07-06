package parser

import (
	"github.com/brandonksides/grundfunken/models"
)

type IdentifierExpression struct {
	name string
	loc  models.SourceLocation
}

func (ie *IdentifierExpression) Type(tb models.TypeBindings) (models.Type, *models.InterpreterError) {
	if _, ok := tb[ie.name]; !ok {
		return nil, &models.InterpreterError{
			Message:        "cannot type unbound identifier",
			SourceLocation: ie.loc,
		}
	}

	return tb[ie.name], nil
}

func (ie *IdentifierExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret, ok := map[string]any(bindings)[ie.name]
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "cannot evaluate unbound identifier",
			SourceLocation: ie.loc,
		}
	}

	return ret, nil
}

func (ie *IdentifierExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}
