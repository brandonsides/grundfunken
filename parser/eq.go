package parser

import (
	"errors"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type EqExpression struct {
	Left  models.Expression
	Right models.Expression
	Op    EqOp
}

type EqOp struct {
	Type           EqOpType
	SourceLocation models.SourceLocation
}

type EqOpType int

const (
	EQ_OP_EQUAL EqOpType = iota
	EQ_OP_NOT_EQUAL
)

func (ee *EqExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := ee.Left.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	v2, err := ee.Right.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	switch ee.Op.Type {
	case EQ_OP_EQUAL:
		return v1 == v2, nil
	case EQ_OP_NOT_EQUAL:
		return v1 != v2, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
			SourceLocation: ee.Op.SourceLocation,
		}
	}
}

func (ee *EqExpression) SourceLocation() models.SourceLocation {
	return ee.Left.SourceLocation()
}

func parseEqExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	left, rest, err := parseCmpExpression(toks)
	if err != nil {
		return nil, toks, err
	}
	if len(rest) == 0 {
		return left, rest, nil
	}

	return foldEq(left, rest)
}

func foldEq(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	var op *EqOp
	op, rest, err = parseEqOp(toks)
	if err != nil {
		return nil, rest, err
	} else if op == nil {
		return first, toks, nil
	}

	right, rest, err := parseEqExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &EqExpression{
		Left:  first,
		Op:    *op,
		Right: right,
	}, rest, nil
}

func parseEqOp(toks []tokens.Token) (op *EqOp, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.IS {
		return nil, toks, nil
	}
	rest = toks[1:]

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type == tokens.NOT {
		op = &EqOp{
			Type:           EQ_OP_NOT_EQUAL,
			SourceLocation: rest[0].SourceLocation,
		}
		rest = rest[1:]
	} else {
		op = &EqOp{
			Type:           EQ_OP_EQUAL,
			SourceLocation: rest[0].SourceLocation,
		}
	}

	return op, rest, nil
}
