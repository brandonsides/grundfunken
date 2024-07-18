package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type AndExpression struct {
	Left  expressions.Expression
	Right expressions.Expression
}

func (ae *AndExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	leftType, err := ae.Left.Type(tb)
	if err != nil {
		return nil, err
	}

	if leftType != types.PrimitiveTypeBool {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator 'and' cannot be applied to type %s", leftType),
			SourceLocation: ae.Left.SourceLocation(),
		}
	}

	rightType, err := ae.Right.Type(tb)
	if err != nil {
		return nil, err
	}

	if rightType != types.PrimitiveTypeBool {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator 'and' cannot be applied to type %s", rightType),
			SourceLocation: ae.Right.SourceLocation(),
		}
	}

	return types.PrimitiveTypeBool, nil
}

func (ae *AndExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	v1, err := ae.Left.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Bool, ok := v1.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v1),
			SourceLocation: ae.Left.SourceLocation(),
		}
	}

	if !v1Bool {
		return false, nil
	}

	v2, err := ae.Right.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Bool, ok := v2.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v2),
			SourceLocation: ae.Right.SourceLocation(),
		}
	}

	return v2Bool, nil
}

func (ae *AndExpression) SourceLocation() *models.SourceLocation {
	return ae.Left.SourceLocation()
}

func parseAndExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	left, err := parseEqExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldAnd(left, toks)
}

func foldAnd(first expressions.Expression, toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.AND {
		return first, nil
	}
	toks.Pop()

	right, err := parseAndExpression(toks)
	if err != nil {
		return nil, err
	}

	return &AndExpression{
		Left:  first,
		Right: right,
	}, nil
}
