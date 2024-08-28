package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type MatchExpression struct {
	On   expressions.Expression
	Arms []MatchArm
	As   string
	loc  *models.SourceLocation
}

type MatchArm struct {
	Type types.Type
	Exp  expressions.Expression
}

func (me *MatchExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	_, err := me.On.Type(tb)
	if err != nil {
		return nil, err
	}

	typs := make([]types.Type, 0, len(me.Arms))
	for _, arm := range me.Arms {
		newTB := make(types.TypeBindings)
		for k, v := range tb {
			newTB[k] = v
		}
		newTB[me.As] = arm.Type

		armType, err := arm.Exp.Type(newTB)
		if err != nil {
			return nil, err
		}

		typs = append(typs, armType)
	}

	return types.Sum(typs...), nil
}

func (me *MatchExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	onVal, err := me.On.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	typ, innerErr := types.TypeOf(onVal)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "cannot determine type of match expression",
			SourceLocation: me.On.SourceLocation(),
			Underlying:     innerErr,
		}
	}

	for _, arm := range me.Arms {
		armSuper, err := types.IsSuperTo(arm.Type, typ)
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "cannot determine type of match expression",
				SourceLocation: arm.Exp.SourceLocation(),
				Underlying:     err,
			}
		}
		if armSuper {
			newBindings := make(expressions.Bindings)
			for k, v := range bindings {
				newBindings[k] = v
			}
			newBindings[me.As] = onVal

			return arm.Exp.Evaluate(bindings)
		}
	}

	return nil, &models.InterpreterError{
		Message: "no match arm found",
	}
}

func (me *MatchExpression) SourceLocation() *models.SourceLocation {
	return me.loc
}

func parseMatchExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected match expression",
			SourceLocation: toks.CurrentSourceLocation(),
			Underlying:     innerErr,
		}
	}

	if tok.Type != tokens.MATCH {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected match expression",
			SourceLocation: &tok.SourceLocation,
		}
	}

	tok, innerErr = toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected match identifier",
			SourceLocation: toks.CurrentSourceLocation(),
			Underlying:     innerErr,
		}
	}

	if tok.Type != tokens.IDENTIFIER {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected identifier",
			SourceLocation: &tok.SourceLocation,
		}
	}
	id := tok.Value

	tok, innerErr = toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected match identifier",
			SourceLocation: toks.CurrentSourceLocation(),
			Underlying:     innerErr,
		}
	}

	if tok.Type != tokens.ON {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected 'on'",
			SourceLocation: &tok.SourceLocation,
		}
	}

	onExp, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	ret := &MatchExpression{
		On:  onExp,
		As:  id,
		loc: onExp.SourceLocation(),
	}
	for {
		tok, ok := toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "expected match arms",
				SourceLocation: onExp.SourceLocation(),
			}
		}

		if tok.Type != tokens.CASE {
			break
		}
		toks.Pop()

		typ, innerErr := parseType(toks)
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected type",
				SourceLocation: &tok.SourceLocation,
				Underlying:     innerErr,
			}
		}

		exp, err := ParseExpression(toks)
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected expression",
				SourceLocation: &tok.SourceLocation,
				Underlying:     err,
			}
		}

		ret.Arms = append(ret.Arms, MatchArm{
			Type: typ,
			Exp:  exp,
		})
	}

	return ret, nil
}
