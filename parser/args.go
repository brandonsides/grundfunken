package parser

import (
	"errors"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseExpressions(toks []tokens.Token) (exps []models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exps = make([]models.Expression, 0)
	var exp models.Expression
	for exp, rest, err = ParseExpression(toks); err == nil; exp, rest, err = ParseExpression(rest) {
		if exp == nil {
			return exps, rest, nil
		}

		exps = append(exps, exp)
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		if rest[0].Type != tokens.COMMA {
			break
		}
		rest = rest[1:]
	}
	if err != nil {
		return nil, rest, err
	}
	return exps, rest, nil
}
