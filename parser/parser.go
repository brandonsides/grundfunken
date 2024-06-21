package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

func ParseExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	return parseOrExpression(toks)
}
