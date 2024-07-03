package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type NotExpression struct {
	Inner models.Expression
	loc   models.SourceLocation
}

func (ne *NotExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v, err := ne.Inner.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	vBool, ok := v.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected bool",
			SourceLocation: ne.Inner.SourceLocation(),
		}
	}

	return !vBool, nil
}

func (ne *NotExpression) SourceLocation() models.SourceLocation {
	return ne.loc
}

func parseNotExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected token",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type == tokens.NOT {
		toks.Pop()
		inner, err := parseNotExpression(toks)
		if err != nil {
			return nil, err
		}

		return &NotExpression{
			Inner: inner,
			loc:   tok.SourceLocation,
		}, nil
	}

	return parseAtomic(toks)
}
