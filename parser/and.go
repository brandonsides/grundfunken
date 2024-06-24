package parser

import (
	"errors"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type AndExpression struct {
	Left  models.Expression
	Right models.Expression
}

func (ae *AndExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := ae.Left.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Bool, ok := v1.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("expected bool"),
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
			Err:            errors.New("expected bool"),
			SourceLocation: ae.Right.SourceLocation(),
		}
	}

	return v2Bool, nil
}

func (ae *AndExpression) SourceLocation() models.SourceLocation {
	return ae.Left.SourceLocation()
}

func parseAndExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	left, rest, err := parseEqExpression(toks)
	if err != nil {
		return nil, toks, err
	}
	if len(rest) == 0 {
		return left, rest, nil
	}

	return foldAnd(left, rest)
}

func foldAnd(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 || toks[0].Type != tokens.AND {
		return first, toks, nil
	}
	rest = toks[1:]

	right, rest, err := parseAndExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &AndExpression{
		Left:  first,
		Right: right,
	}, rest, nil
}
