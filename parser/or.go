package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type OrExpression struct {
	Left  models.Expression
	Right models.Expression
}

func (oe *OrExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := oe.Left.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Bool, ok := v1.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v1),
			SourceLocation: oe.Left.SourceLocation(),
		}
	}

	if v1Bool {
		return true, nil
	}

	v2, err := oe.Right.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Bool, ok := v2.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v2),
			SourceLocation: oe.Right.SourceLocation(),
		}
	}

	return v2Bool, nil
}

func (oe *OrExpression) SourceLocation() models.SourceLocation {
	return oe.Left.SourceLocation()
}

func parseOrExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	left, rest, err := parseAndExpression(toks)
	if err != nil {
		return nil, toks, err
	}
	if len(rest) == 0 {
		return left, rest, nil
	}

	return foldOr(left, rest)
}

func foldOr(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 || toks[0].Type != tokens.OR {
		return first, toks, nil
	}
	rest = toks[1:]

	right, rest, err := parseOrExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &OrExpression{
		Left:  first,
		Right: right,
	}, rest, nil
}
