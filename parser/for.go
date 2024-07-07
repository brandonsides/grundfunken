package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type ForExpression struct {
	ForClause  models.Expression
	Identifier string
	InClause   models.Expression
	loc        models.SourceLocation
}

func (fe *ForExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	inType, err := fe.InClause.Type(tb)
	if err != nil {
		return nil, err
	}

	inTypeList, ok := inType.(types.ListType)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("for expression in clause must evaluate to a list; got %s", inType),
			SourceLocation: fe.InClause.SourceLocation(),
		}
	}

	innerTB := make(types.TypeBindings)
	for k, v := range tb {
		innerTB[k] = v
	}
	innerTB[fe.Identifier] = inTypeList.ElementType

	forType, err := fe.ForClause.Type(innerTB)
	if err != nil {
		return nil, err
	}

	return types.List(forType), nil
}

func (fe *ForExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	ret := make([]any, 0)

	innerBindings := make(models.Bindings)
	for k, v := range bindings {
		innerBindings[k] = v
	}

	iterableExp, err := fe.InClause.Evaluate(innerBindings)
	if err != nil {
		return nil, err
	}

	iterableExpArr, ok := iterableExp.([]any)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("for expression in clause must evaluate to an array; got %v", iterableExp),
			SourceLocation: fe.InClause.SourceLocation(),
		}
	}

	for _, v := range iterableExpArr {
		innerBindings[fe.Identifier] = v
		retVal, err := fe.ForClause.Evaluate(innerBindings)
		if err != nil {
			return nil, err
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func (fe *ForExpression) SourceLocation() models.SourceLocation {
	return fe.loc
}

func parseForExpression(exp1 models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return exp1, nil
	}
	if tok.Type != tokens.FOR {
		return exp1, nil
	}
	toks.Pop()

	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected for clause",
			Underlying:     innerErr,
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}
	identifier := tok.Value

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected in clause",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.IN {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected in clause",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	exp2, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	return &ForExpression{
		ForClause:  exp1,
		Identifier: identifier,
		InClause:   exp2,
		loc:        beginLoc,
	}, nil
}
