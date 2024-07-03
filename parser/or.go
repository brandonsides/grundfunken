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

func parseOrExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	left, err := parseAndExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldOr(left, toks)
}

func foldOr(first models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.OR {
		return first, nil
	}
	toks.Pop()

	right, err := parseOrExpression(toks)
	if err != nil {
		return nil, err
	}

	return &OrExpression{
		Left:  first,
		Right: right,
	}, nil
}
