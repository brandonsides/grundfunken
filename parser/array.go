package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type ArrayLiteralExpression struct {
	elemType types.Type
	val      []expressions.Expression
	loc      *models.SourceLocation
}

func (ale *ArrayLiteralExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	if ale.elemType == nil {
		ale.elemType = types.PrimitiveTypeAny
	}

	for _, v := range ale.val {
		t, err := v.Type(tb)
		if err != nil {
			return nil, err
		}

		aleSuper, innerErr := types.IsSuperTo(ale.elemType, t)
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "inconsistent array element types",
				SourceLocation: v.SourceLocation(),
				Underlying:     innerErr,
			}
		}
		if !aleSuper {
			return nil, &models.InterpreterError{
				Message:        "inconsistent array element types",
				SourceLocation: v.SourceLocation(),
			}
		}
	}

	return types.List(ale.elemType), nil
}

func (ale *ArrayLiteralExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	ret := make([]any, 0)
	for _, v := range ale.val {
		retVal, err := v.Evaluate(bindings)
		if err != nil {
			return nil, err
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func (ale *ArrayLiteralExpression) SourceLocation() *models.SourceLocation {
	return ale.loc
}

func parseArrayLiteral(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	beginSourceLocation := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "unexpected end of input",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_SQUARE_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: &tok.SourceLocation,
		}
	}
	toks.Pop()

	exps, err := parseExpressions(toks)
	if err != nil {
		return nil, err
	}

	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "to terminate array literal",
			SourceLocation: beginSourceLocation,
			Underlying: &models.InterpreterError{
				Message:        "expected closing square bracket",
				Underlying:     innerErr,
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}

	if tok.Type != tokens.RIGHT_SQUARE_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "to terminate array literal",
			SourceLocation: beginSourceLocation,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected closing square bracket",
				SourceLocation: &tok.SourceLocation,
			},
		}
	}

	typ, innerErr := parseType(toks)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "failed to parse array type",
			SourceLocation: beginSourceLocation,
			Underlying:     innerErr,
		}
	}

	exp = &ArrayLiteralExpression{
		val:      exps,
		loc:      beginSourceLocation,
		elemType: typ,
	}

	return exp, nil
}
