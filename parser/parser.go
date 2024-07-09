package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func ParseExpression(toks *tokens.TokenStack) (lines map[string]string, exp models.Expression, err *models.InterpreterError) {
	return parseOrExpression(toks)
}
