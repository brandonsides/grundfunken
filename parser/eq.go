package parser

import (
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
			Message:        "invalid operator",
			SourceLocation: ee.Op.SourceLocation,
		}
	}
}

func (ee *EqExpression) SourceLocation() models.SourceLocation {
	return ee.Left.SourceLocation()
}

func parseEqExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	left, err := parseCmpExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldEq(left, toks)
}

func foldEq(first models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	var op *EqOp
	op, err = parseEqOp(toks)
	if err != nil {
		return nil, err
	} else if op == nil {
		return first, nil
	}

	right, err := parseEqExpression(toks)
	if err != nil {
		return nil, err
	}

	return &EqExpression{
		Left:  first,
		Op:    *op,
		Right: right,
	}, nil
}

func parseEqOp(toks *tokens.TokenStack) (op *EqOp, err *models.InterpreterError) {
	tok := toks.Peek()
	if tok == nil || tok.Type != tokens.IS {
		return nil, nil
	}
	toks.Pop()

	tok = toks.Pop()
	if tok == nil {
		return nil, &models.InterpreterError{
			Message:        "expected expression",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}
	if tok.Type == tokens.NOT {
		op = &EqOp{
			Type:           EQ_OP_NOT_EQUAL,
			SourceLocation: tok.SourceLocation,
		}
	} else {
		op = &EqOp{
			Type:           EQ_OP_EQUAL,
			SourceLocation: tok.SourceLocation,
		}
	}

	return op, nil
}
