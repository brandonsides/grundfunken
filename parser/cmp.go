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

func parseCmpExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	exp, rest, err = parseAddExpression(toks)
	if err != nil {
		return nil, rest, err
	}

	return foldCmp(exp, rest)
}

func foldCmp(first models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return first, toks, nil
	}

	op := CmpOp{
		SourceLocation: toks[0].SourceLocation,
	}

	switch toks[0].Type {
	case tokens.LEFT_ANGLE_BRACKET:
		if len(toks) == 1 {
			return nil, toks, &models.InterpreterError{
				Message: "unexpected end of input",
			}
		}
		next := toks[1]
		if next.Type == tokens.EQUAL {
			op.Type = CMP_OP_TYPE_LESS_EQUAL
			rest = toks[2:]
		} else {
			op.Type = CMP_OP_TYPE_LESS
			rest = toks[1:]
		}
	case tokens.RIGHT_ANGLE_BRACKET:
		if len(toks) == 1 {
			return nil, toks, &models.InterpreterError{
				Message: "unexpected end of input",
			}
		}
		next := toks[1]
		if next.Type == tokens.EQUAL {
			op.Type = CMP_OP_GREATER_EQUAL
			rest = toks[2:]
		} else {
			op.Type = CMP_OP_GREATER
			rest = toks[1:]
		}
	default:
		return first, toks, nil
	}

	exp2, rest, err := parseCmpExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &CmpExpression{
		first:  first,
		op:     op,
		second: exp2,
	}, rest, nil
}
