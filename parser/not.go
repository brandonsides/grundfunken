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

func parseNotExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Message: "expected token",
		}
	}

	if toks[0].Type == tokens.NOT {
		inner, rest, err := parseNotExpression(toks[1:])
		if err != nil {
			return nil, toks, err
		}

		return &NotExpression{
			Inner: inner,
			loc:   toks[0].SourceLocation,
		}, rest, nil
	}

	return parseAtomic(toks)
}
