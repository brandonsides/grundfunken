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

func (fce *FunctionCallExpression) Type(tb models.TypeBindings) (models.Type, *models.InterpreterError) {
	funType, err := fce.Function.Type(tb)
	if err != nil {
		return nil, err
	}

	if funType != models.PrimitiveTypeFunction {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("cannot call non-function %s", funType.String()),
			SourceLocation: fce.Function.SourceLocation(),
		}
	}

	for _, arg := range fce.Args {
		_, err := arg.Type(tb)
		if err != nil {
			return nil, err
		}
	}

	return models.PrimitiveTypeAny, nil
}

func (fce *FunctionCallExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	f, err := fce.Function.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	fun, ok := f.(models.Function)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("cannot call non-function %v", f),
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
