package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

type MulExpression struct {
	op     tokens.Token
	first  models.Expression
	second models.Expression
}

func (me *MulExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := me.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Muller, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("operator '*' cannot be applied to first operand"),
			SourceLocation: me.first.SourceLocation(),
		}
	}

	v2, err := me.second.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Muller, ok := v2.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("operator '*' cannot be applied to second operand"),
			SourceLocation: me.second.SourceLocation(),
		}
	}

	switch me.op.Type {
	case tokens.STAR:
		return v1Muller * v2Muller, nil
	case tokens.SLASH:
		return v1Muller / v2Muller, nil
	case tokens.PERCENT:
		return v1Muller % v2Muller, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
			SourceLocation: me.op.SourceLocation,
		}
	}
}

func (me *MulExpression) SourceLocation() models.SourceLocation {
	return me.first.SourceLocation()
}

func parseMulExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exp, rest, err = parseNotExpression(toks)
	if err != nil {
		return nil, rest, err
	}

	return foldMul(exp, rest)
}

func foldMul(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	if toks[0].Type != tokens.STAR && toks[0].Type != tokens.SLASH && toks[0].Type != tokens.PERCENT {
		return first, toks, nil
	}

	rest = toks[1:]

	right, rest, err := parseMulExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &MulExpression{
		first:  first,
		op:     toks[0],
		second: right,
	}, rest, nil
}
