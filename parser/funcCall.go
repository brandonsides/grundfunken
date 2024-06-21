package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
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
			Err:            errors.New("cannot call non-function"),
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

	return fun.Call(argArray)
}

func (fce *FunctionCallExpression) SourceLocation() models.SourceLocation {
	return fce.loc
}
