package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type ForExpression struct {
	ForClause  models.Expression
	Identifier string
	InClause   models.Expression
	loc        models.SourceLocation
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
			Message:        "for expression must evaluate to an array",
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

func parseForExpression(exp1 models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return exp1, toks, nil
	}

	if toks[0].Type != tokens.FOR {
		return exp1, toks, nil
	}

	rest = toks[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	identifier := rest[0].Value
	rest = rest[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	if rest[0].Type != tokens.IN {
		return nil, rest, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &ForExpression{
		ForClause:  exp1,
		Identifier: identifier,
		InClause:   exp2,
		loc:        toks[0].SourceLocation,
	}, rest, nil
}
