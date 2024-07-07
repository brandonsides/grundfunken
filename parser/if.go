package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type IfExpression struct {
	Condition models.Expression
	Then      models.Expression
	Else      models.Expression
	loc       models.SourceLocation
}

func (ie *IfExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	condType, err := ie.Condition.Type(tb)
	if err != nil {
		return nil, err
	}

	if condType != types.PrimitiveTypeBool {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("if condition must evaluate to a boolean; got %s", condType),
			SourceLocation: ie.Condition.SourceLocation(),
		}
	}

	thenType, err := ie.Then.Type(tb)
	if err != nil {
		return nil, err
	}

	elseType, err := ie.Else.Type(tb)
	if err != nil {
		return nil, err
	}

	return types.Sum(thenType, elseType), nil
}

func (ie *IfExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	cond, err := ie.Condition.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	condBool, ok := cond.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("if condition must evaluate to a boolean; got %v", cond),
			SourceLocation: ie.Condition.SourceLocation(),
		}
	}

	if condBool {
		return ie.Then.Evaluate(bindings)
	}

	return ie.Else.Evaluate(bindings)
}

func (ie *IfExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}

func parseIfExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected \"if\" expression",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.IF {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected \"if\" clause",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	exp1, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "in \"if\" expression",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "expected \"then\" clause",
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}

	if tok.Type != tokens.THEN {
		return nil, &models.InterpreterError{
			Message:        "in \"if\" expression",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected \"then\" clause",
				SourceLocation: tok.SourceLocation,
			},
		}
	}
	toks.Pop()

	exp2, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "in \"if\" expression",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "expected \"else\" clause",
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}

	if tok.Type != tokens.ELSE {
		return nil, &models.InterpreterError{
			Message:        "in \"if\" expression",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected \"else\" clause",
				SourceLocation: tok.SourceLocation,
			},
		}
	}
	toks.Pop()

	exp3, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	return &IfExpression{
		Condition: exp1,
		Then:      exp2,
		Else:      exp3,
		loc:       beginLoc,
	}, nil
}
