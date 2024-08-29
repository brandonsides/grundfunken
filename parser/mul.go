package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type MulExpression struct {
	op     tokens.Token
	first  expressions.Expression
	second expressions.Expression
}

func (me *MulExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	firstType, err := me.first.Type(tb)
	if err != nil {
		return nil, err
	}

	if firstType != types.PrimitiveTypeInt {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator '%s' cannot be applied to type %s", me.op.Value, firstType),
			SourceLocation: me.first.SourceLocation(),
		}
	}

	secondType, err := me.second.Type(tb)
	if err != nil {
		return nil, err
	}

	if secondType != types.PrimitiveTypeInt {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator '%s' cannot be applied to type %s", me.op.Value, secondType),
			SourceLocation: me.second.SourceLocation(),
		}
	}

	return types.PrimitiveTypeInt, nil
}

func (me *MulExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	v1, err := me.first.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Muller, ok := v1.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator '%s' cannot be applied to first operand %v", me.op.Value, v1),
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
			Message:        fmt.Sprintf("operator '%s' cannot be applied to second operand %v", me.op.Value, v2),
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
			Message:        "invalid operator",
			SourceLocation: &me.op.SourceLocation,
		}
	}
}

func (me *MulExpression) SourceLocation() *models.SourceLocation {
	return me.first.SourceLocation()
}

func parseMulExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	exp, err = parseNotExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldMul(exp, toks)
}

func foldMul(first expressions.Expression, toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.STAR && tok.Type != tokens.SLASH && tok.Type != tokens.PERCENT {
		return first, nil
	}
	if first == nil {
		return nil, &models.InterpreterError{
			Message:        "expected expression",
			SourceLocation: &tok.SourceLocation,
		}
	}

	toks.Pop()

	next, err := parseNotExpression(toks)
	if err != nil {
		return first, err
	}

	withNext := &MulExpression{
		op: tok,
		first: first,
		second: next,
	}

	return foldMul(withNext, toks)
}
