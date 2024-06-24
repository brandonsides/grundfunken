package parser

import (
	"errors"

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

func parseArrayLiteral(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.LEFT_SQUARE_BRACKET {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	exps, rest, err := parseExpressions(rest)
	if err != nil {
		return nil, rest, err
	}
	if rest[0].Type != tokens.RIGHT_SQUARE_BRACKET {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}
	rest = rest[1:]
	exp = &ArrayLiteralExpression{
		val: exps,
		loc: toks[0].SourceLocation,
	}

	return exp, rest, nil
}
