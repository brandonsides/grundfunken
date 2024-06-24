package parser

import (
	"errors"
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type AddExpression struct {
	op     tokens.Token
	first  models.Expression
	second models.Expression
}

func (ae *AddExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := ae.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Adder, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            fmt.Errorf("operator '%s' cannot be applied to first operand", ae.op.Value),
			SourceLocation: ae.first.SourceLocation(),
		}
	}

	v2, err := ae.second.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Adder, ok := v2.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            fmt.Errorf("operator '%s' cannot be applied to second operand", ae.op.Value),
			SourceLocation: ae.second.SourceLocation(),
		}
	}

	switch ae.op.Type {
	case tokens.PLUS:
		return v1Adder + v2Adder, nil
	case tokens.MINUS:
		return v1Adder - v2Adder, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
			SourceLocation: ae.op.SourceLocation,
		}
	}
}

func (ae *AddExpression) SourceLocation() models.SourceLocation {
	return ae.first.SourceLocation()
}

func parseAddExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exp, rest, err = parseMulExpression(toks)
	if err != nil {
		return nil, rest, err
	}

	return foldAdd(exp, rest)
}

func foldAdd(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	if toks[0].Type != tokens.PLUS && toks[0].Type != tokens.MINUS {
		return first, toks, nil
	}
	op := toks[0]

	rest = toks[1:]

	var withNext models.Expression
	var next models.Expression
	next, rest, err = parseMulExpression(rest)
	if err != nil {
		return first, rest, err
	}

	withNext = &AddExpression{
		op:     op,
		first:  first,
		second: next,
	}

	return foldAdd(withNext, rest)
}
