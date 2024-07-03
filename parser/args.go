package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseExpressions(toks *tokens.TokenStack) (exps []models.Expression, err *models.InterpreterError) {
	exps = make([]models.Expression, 0)
	var exp models.Expression
	for exp, err = ParseExpression(toks); err == nil; exp, err = ParseExpression(toks) {
		if exp == nil {
			return exps, nil
		}

		exps = append(exps, exp)
		tok, ok := toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "after expression in expression list",
				Underlying:     err,
				SourceLocation: exp.SourceLocation(),
			}
		}

		if tok.Type != tokens.COMMA {
			break
		}

		toks.Pop()
	}
	if err != nil {
		return nil, err
	}
	return exps, nil
}
