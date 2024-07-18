package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
)

type FunctionCallExpression struct {
	Function expressions.Expression
	Args     []expressions.Expression
	loc      *models.SourceLocation
}

func (fce *FunctionCallExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	targetType, err := fce.Function.Type(tb)
	if err != nil {
		return nil, err
	}

	funType, ok := targetType.(types.FuncType)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("cannot call non-function %s", targetType.String()),
			SourceLocation: fce.Function.SourceLocation(),
		}
	}

	if len(fce.Args) != len(funType.ArgTypes) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected %d arguments, got %d", len(funType.ArgTypes), len(fce.Args)),
			SourceLocation: fce.SourceLocation(),
		}
	}

	for i, arg := range fce.Args {
		t, err := arg.Type(tb)
		if err != nil {
			return nil, err
		}

		funSuper, innerErr := types.IsSuperTo(funType.ArgTypes[i], t)
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        fmt.Sprintf("expected %s, got %s", funType.ArgTypes[i].String(), t.String()),
				SourceLocation: arg.SourceLocation(),
				Underlying:     innerErr,
			}
		}
		if !funSuper {
			return nil, &models.InterpreterError{
				Message:        fmt.Sprintf("expected %s, got %s", funType.ArgTypes[i].String(), t.String()),
				SourceLocation: arg.SourceLocation(),
			}
		}
	}

	return funType.ReturnType, nil
}

func (fce *FunctionCallExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	f, err := fce.Function.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	fun, ok := f.(types.Function)
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

func (fce *FunctionCallExpression) SourceLocation() *models.SourceLocation {
	return fce.loc
}
