package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
)

type FunctionCallExpression struct {
	Function models.Expression
	Args     []models.Expression
	loc      models.SourceLocation
}

func (fce *FunctionCallExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	f, err := fce.Function.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	fun, ok := f.(models.Function)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "cannot call non-function",
			SourceLocation: fce.Function.SourceLocation(),
		}
	}

	argArray := make([]any, len(fce.Args))
	for i, arg := range fce.Args {
		val, err := arg.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		argArray[i] = val
	}

	ret, innerErr := fun.Call(argArray)
	if innerErr != nil {
		msg := "in call to anonymous function"
		if identifierExpression, ok := fce.Function.(*IdentifierExpression); ok {
			msg = fmt.Sprintf("in call to function \"%s\"", identifierExpression.name)
		}
		return nil, &models.InterpreterError{
			Message:        msg,
			Underlying:     innerErr,
			SourceLocation: fce.SourceLocation(),
		}
	}
	return ret, nil
}

func (fce *FunctionCallExpression) SourceLocation() models.SourceLocation {
	return fce.loc
}
