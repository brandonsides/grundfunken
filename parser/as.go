package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
)

type AsExpression struct {
	exp   expressions.Expression
	asLoc models.SourceLocation
	typ   types.Type
}

func (ae *AsExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	ulTyp, err := ae.exp.Type(tb)
	if err != nil {
		return nil, err
	}

	canCast, innerErr := types.IsSuperTo(ulTyp, ae.typ)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Underlying: &models.InterpreterError{
				Message:    fmt.Sprintf("failed to determine if %v is a supertype of %v", ulTyp, ae.typ),
				Underlying: innerErr,
			},
			SourceLocation: &ae.asLoc,
			Message:        "in \"as\" expression",
		}
	}

	if !canCast {
		return nil, &models.InterpreterError{
			Message:        "in \"as\" expression",
			Underlying:     fmt.Errorf("expression of type %v can never be of asserted type %v", ulTyp, ae.typ),
			SourceLocation: &ae.asLoc,
		}
	}

	return ae.typ, nil
}

func (ae *AsExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	ret, err := ae.exp.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	retTyp, innerErr := types.TypeOf(ret)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message: "in \"as\" expression",
			Underlying: &models.InterpreterError{
				Message:        fmt.Sprintf("failed to determine runtime type of value %v"),
				Underlying:     innerErr,
				SourceLocation: ae.SourceLocation(),
			},
			SourceLocation: &ae.asLoc,
		}
	}

	canCast, innerErr := types.IsSuperTo(ae.typ, retTyp)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message: "in \"as\" expression",
			Underlying: &models.InterpreterError{
				Message:        fmt.Sprintf("failed to determine if %v is a supertype of %v", ae.typ, retTyp),
				Underlying:     innerErr,
				SourceLocation: ae.SourceLocation(),
			},
			SourceLocation: &ae.asLoc,
		}
	}

	if !canCast {
		return nil, &models.InterpreterError{
			Message: "in \"as\" expression",
			Underlying: &models.InterpreterError{
				Message:        fmt.Sprintf("%v is not of assumed type %v", ret, ae.typ),
				SourceLocation: ae.exp.SourceLocation(),
			},
			SourceLocation: &ae.asLoc,
		}
	}

	return ret, nil
}

func (ae *AsExpression) SourceLocation() *models.SourceLocation {
	return ae.exp.SourceLocation()
}
