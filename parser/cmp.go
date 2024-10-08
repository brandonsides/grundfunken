package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type CmpExpression struct {
	first  expressions.Expression
	op     CmpOp
	second expressions.Expression
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

func (ce *CmpOpType) String() string {
	switch *ce {
	case CMP_OP_TYPE_LESS:
		return "<"
	case CMP_OP_TYPE_LESS_EQUAL:
		return "<="
	case CMP_OP_GREATER_EQUAL:
		return ">="
	case CMP_OP_GREATER:
		return ">"
	default:
		return "unknown"
	}
}

func (ce *CmpExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	firstType, err := ce.first.Type(tb)
	if err != nil {
		return nil, err
	}

	if firstType != types.PrimitiveTypeInt {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator '%s' cannot be applied to type %s", ce.op.Type.String(), firstType.String()),
			SourceLocation: ce.first.SourceLocation(),
		}
	}

	secondType, err := ce.second.Type(tb)
	if err != nil {
		return nil, err
	}

	if secondType != types.PrimitiveTypeInt {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator '%s' cannot be applied to type %s", ce.op.Type.String(), secondType.String()),
			SourceLocation: ce.second.SourceLocation(),
		}
	}

	return types.PrimitiveTypeBool, nil
}

func (ce *CmpExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
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
			SourceLocation: &ce.op.SourceLocation,
		}
	}
}

func (ce *CmpExpression) SourceLocation() *models.SourceLocation {
	return ce.first.SourceLocation()
}

func parseCmpExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	exp, err = parseAddExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldCmp(exp, toks)
}

func foldCmp(first expressions.Expression, toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()
	tok, ok := toks.Peek()
	if !ok {
		return first, nil
	}

	op := CmpOp{
		SourceLocation: tok.SourceLocation,
	}

	switch tok.Type {
	case tokens.LEFT_ANGLE_BRACKET:
		toks.Pop()
		tok, ok = toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "after comparison operator",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "expected expression",
					SourceLocation: toks.CurrentSourceLocation(),
				},
			}
		}
		if tok.Type == tokens.EQUAL {
			op.Type = CMP_OP_TYPE_LESS_EQUAL
			toks.Pop()
		} else {
			op.Type = CMP_OP_TYPE_LESS
		}
	case tokens.RIGHT_ANGLE_BRACKET:
		toks.Pop()
		tok, ok = toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "after comparison operator",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "expected expression",
					SourceLocation: toks.CurrentSourceLocation(),
				},
			}
		}
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
