package parser

import (
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
			Message:        fmt.Sprintf("operator '%s' cannot be applied to first operand %v", ae.op.Value, v1),
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
			Message:        fmt.Sprintf("operator '%s' cannot be applied to second operand %v", ae.op.Value, v2),
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
			Message:        "invalid operator",
			SourceLocation: ae.op.SourceLocation,
		}
	}
}

func (ae *AddExpression) SourceLocation() models.SourceLocation {
	return ae.first.SourceLocation()
}

func parseAddExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	exp, err = parseMulExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldAdd(exp, toks)
}

func foldAdd(first models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.PLUS && tok.Type != tokens.MINUS {
		return first, nil
	}
	toks.Pop()

	var withNext models.Expression
	var next models.Expression
	next, err = parseMulExpression(toks)
	if err != nil {
		return first, err
	}

	withNext = &AddExpression{
		op:     tok,
		first:  first,
		second: next,
	}

	return foldAdd(withNext, toks)
}
