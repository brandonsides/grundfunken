package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type CmpExpression struct {
	first  models.Expression
	op     CmpOp
	second models.Expression
}

type CmpOp struct {
	Type           CmpOpType
	SourceLocation models.SourceLocation
}

type CmpOpType int

const (
	CMP_OP_TYPE_LESS CmpOpType = iota
	CMP_OP_TYPE_LESS_EQUAL
	CMP_OP_GREATER_EQUAL
	CMP_OP_GREATER
)

func (ce *CmpExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := ce.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	v2, err := ce.second.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	v1Int, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected int; got %v", v1),
			SourceLocation: ce.first.SourceLocation(),
		}
	}

	v2Int, ok := v2.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected int; got %v", v2),
			SourceLocation: ce.second.SourceLocation(),
		}
	}

	switch ce.op.Type {
	case CMP_OP_TYPE_LESS:
		return v1Int < v2Int, nil
	case CMP_OP_TYPE_LESS_EQUAL:
		return v1Int <= v2Int, nil
	case CMP_OP_GREATER_EQUAL:
		return v1Int >= v2Int, nil
	case CMP_OP_GREATER:
		return v1Int > v2Int, nil
	default:
		return nil, &models.InterpreterError{
			Message:        "invalid operator",
			SourceLocation: ce.op.SourceLocation,
		}
	}
}

func (ce *CmpExpression) SourceLocation() models.SourceLocation {
	return ce.first.SourceLocation()
}

func parseCmpExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	exp, err = parseAddExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldCmp(exp, toks)
}

func foldCmp(first models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok := toks.Peek()
	if tok == nil {
		return first, nil
	}
	toks.Pop()

	op := CmpOp{
		SourceLocation: tok.SourceLocation,
	}

	switch tok.Type {
	case tokens.LEFT_ANGLE_BRACKET:
		tok = toks.Peek()
		if tok.Type == tokens.EQUAL {
			op.Type = CMP_OP_TYPE_LESS_EQUAL
			toks.Pop()
		} else {
			op.Type = CMP_OP_TYPE_LESS
		}
	case tokens.RIGHT_ANGLE_BRACKET:
		tok = toks.Peek()
		if tok.Type == tokens.EQUAL {
			op.Type = CMP_OP_GREATER_EQUAL
			toks.Pop()
		} else {
			op.Type = CMP_OP_GREATER
		}
	default:
		return first, nil
	}

	exp2, err := parseCmpExpression(toks)
	if err != nil {
		return nil, err
	}

	return &CmpExpression{
		first:  first,
		op:     op,
		second: exp2,
	}, nil
}
