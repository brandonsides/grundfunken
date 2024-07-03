package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type ArrayLiteralExpression struct {
	val []models.Expression
	loc models.SourceLocation
}

func (ale *ArrayLiteralExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret := make([]any, 0)
	for _, v := range ale.val {
		retVal, err := v.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func (ale *ArrayLiteralExpression) SourceLocation() models.SourceLocation {
	return ale.loc
}

func parseArrayLiteral(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginSourceLocation := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "unexpected end of input",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_SQUARE_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	exps, err := parseExpressions(toks)
	if err != nil {
		return nil, err
	}

	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "to terminate array literal",
			SourceLocation: beginSourceLocation,
			Underlying: &models.InterpreterError{
				Message:        "expected closing square bracket",
				Underlying:     innerErr,
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}

	if tok.Type != tokens.RIGHT_SQUARE_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "to terminate array literal",
			SourceLocation: beginSourceLocation,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected closing square bracket",
				SourceLocation: tok.SourceLocation,
			},
		}
	}

	exp = &ArrayLiteralExpression{
		val: exps,
		loc: beginSourceLocation,
	}

	return exp, nil
}
