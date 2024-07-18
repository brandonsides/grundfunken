package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
)

type IdentifierExpression struct {
	name string
	loc  models.SourceLocation
}

func (ie *IdentifierExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	if _, ok := tb[ie.name]; !ok {
		return nil, &models.InterpreterError{
			Message:        "cannot type unbound identifier",
			SourceLocation: &ie.loc,
		}
	}

	return tb[ie.name], nil
}

func (ie *IdentifierExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	ret, ok := map[string]any(bindings)[ie.name]
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "cannot evaluate unbound identifier",
			SourceLocation: &ie.loc,
		}
	}

	return ret, nil
}

func (ie *IdentifierExpression) SourceLocation() *models.SourceLocation {
	return &ie.loc
}
