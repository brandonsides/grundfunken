package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
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

	switch ce.op.Type {
	case CMP_OP_TYPE_LESS:
		v1Int, ok := v1.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.first.SourceLocation(),
			}
		}

		v2Int, ok := v2.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.second.SourceLocation(),
			}
		}

		return v1Int < v2Int, nil
	case CMP_OP_TYPE_LESS_EQUAL:
		v1Int, ok := v1.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.first.SourceLocation(),
			}
		}

		v2Int, ok := v2.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.second.SourceLocation(),
			}
		}

		return v1Int <= v2Int, nil
	case CMP_OP_GREATER_EQUAL:
		v1Int, ok := v1.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.first.SourceLocation(),
			}
		}

		v2Int, ok := v2.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.second.SourceLocation(),
			}
		}

		return v1Int >= v2Int, nil
	case CMP_OP_GREATER:
		v1Int, ok := v1.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.first.SourceLocation(),
			}
		}

		v2Int, ok := v2.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: ce.second.SourceLocation(),
			}
		}

		return v1Int > v2Int, nil
	default:
		return nil, &models.InterpreterError{
			Err:            errors.New("invalid operator"),
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
				Err: errors.New("unexpected end of input"),
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
				Err: errors.New("unexpected end of input"),
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
