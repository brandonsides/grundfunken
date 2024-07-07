package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type OrExpression struct {
	Left  models.Expression
	Right models.Expression
}

func (oe *OrExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	leftType, err := oe.Left.Type(tb)
	if err != nil {
		return nil, err
	}

	if leftType != types.PrimitiveTypeBool {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator 'or' cannot be applied to type %s", leftType),
			SourceLocation: oe.Left.SourceLocation(),
		}
	}

	rightType, err := oe.Right.Type(tb)
	if err != nil {
		return nil, err
	}

	if rightType != types.PrimitiveTypeBool {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("operator 'or' cannot be applied to type %s", rightType),
			SourceLocation: oe.Right.SourceLocation(),
		}
	}

	return types.PrimitiveTypeBool, nil
}

func (oe *OrExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	v1, err := oe.Left.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v1Bool, ok := v1.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v1),
			SourceLocation: oe.Left.SourceLocation(),
		}
	}

	if v1Bool {
		return true, nil
	}

	v2, err := oe.Right.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	v2Bool, ok := v2.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected bool; got %v", v2),
			SourceLocation: oe.Right.SourceLocation(),
		}
	}

	return v2Bool, nil
}

func (oe *OrExpression) SourceLocation() models.SourceLocation {
	return oe.Left.SourceLocation()
}

func parseOrExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	left, err := parseAndExpression(toks)
	if err != nil {
		return nil, err
	}

	return foldOr(left, toks)
}

func foldOr(first models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok, ok := toks.Peek()
	if !ok || tok.Type != tokens.OR {
		return first, nil
	}
	if first == nil {
		return nil, &models.InterpreterError{
			Message:        "expected expression",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	right, err := parseOrExpression(toks)
	if err != nil {
		return nil, err
	}

	return &OrExpression{
		Left:  first,
		Right: right,
	}, nil
}
