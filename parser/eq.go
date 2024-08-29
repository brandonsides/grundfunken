package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type EqExpression struct {
	Left  expressions.Expression
	Right expressions.Expression
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

func (ee *EqExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	_, err := ee.Left.Type(tb)
	if err != nil {
		return nil, err
	}

	_, err = ee.Right.Type(tb)
	if err != nil {
		return nil, err
	}

	return types.PrimitiveTypeBool, nil
}

func (ee *EqExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
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
			SourceLocation: &ee.Op.SourceLocation,
		}
	}
}

func (ee *EqExpression) SourceLocation() *models.SourceLocation {
	return ee.Left.SourceLocation()
}

func parseEqExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	left, err := parseCmpExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldEq(left, toks)
}

func foldEq(first expressions.Expression, toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	var op *EqOp
	op, err = parseEqOp(toks)
	if err != nil {
		return nil, err
	} else if op == nil {
		return first, nil
	}

	next, err := parseCmpExpression(toks)
	if err != nil {
		return first, err
	}

	withNext := &EqExpression{
		Op:    *op,
		Left:  first,
		Right: next,
	}

	return foldEq(withNext, toks)
}

func parseEqOp(toks *tokens.TokenStack) (op *EqOp, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.IS {
		return nil, nil
	}
	toks.Pop()

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "after \"is\" operator",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "expected expression",
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}
	if tok.Type == tokens.NOT {
		toks.Pop()
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
