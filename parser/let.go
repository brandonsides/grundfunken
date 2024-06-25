package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type BindingExpression struct {
	Identifier string
	Expression models.Expression
}

type LetExpression struct {
	loc        models.SourceLocation
	LetClauses []BindingExpression
	InClause   models.Expression
}

func (le *LetExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	newBindings := make(models.Bindings)
	for k, v := range bindings {
		newBindings[k] = v
	}

	for _, bindingExp := range le.LetClauses {
		k, v := bindingExp.Identifier, bindingExp.Expression
		val, err := v.Evaluate(newBindings)
		if err != nil {
			return nil, err
		}

		newBindings[k] = val

		if funcVal, ok := val.(*FuncValue); ok {
			funcVal.Bindings[k] = val
		}
	}

	return le.InClause.Evaluate(newBindings)
}

func (le *LetExpression) SourceLocation() models.SourceLocation {
	return le.loc
}

func parseLetExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) < 1 {
		return nil, toks, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	if toks[0].Type != tokens.LET {
		return nil, toks, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: toks[0].SourceLocation,
		}
	}

	if len(rest) == 1 {
		return nil, toks, &models.InterpreterError{
			Message:        "expected identifier after \"let\"",
			SourceLocation: toks[0].SourceLocation,
		}
	}
	rest = toks[1:]

	bindingExpressions := make([]BindingExpression, 0)
	for len(rest) > 0 {
		if rest[0].Type != tokens.IDENTIFIER {
			return nil, rest, &models.InterpreterError{
				Message:        "unexpected token",
				SourceLocation: rest[0].SourceLocation,
			}
		}

		identifier := rest[0].Value
		rest = rest[1:]

		if rest[0].Type != tokens.EQUAL {
			return nil, rest, &models.InterpreterError{
				Message:        "unexpected token",
				SourceLocation: rest[0].SourceLocation,
			}
		}

		if len(rest) == 1 {
			return nil, rest, &models.InterpreterError{
				Message:        fmt.Sprintf("expected expression after \"%s\"", rest[0].Value),
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]

		var exp1 models.Expression
		exp1, rest, err = ParseExpression(rest)
		if err != nil {
			return nil, rest, err
		}

		bindingExpressions = append(bindingExpressions, BindingExpression{
			Identifier: identifier,
			Expression: exp1,
		})

		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Message:        "expected in clause or binding clause after binding expression",
				SourceLocation: exp1.SourceLocation(),
			}
		}

		if rest[0].Type != tokens.COMMA {
			break
		}

		if len(rest) == 1 {
			return nil, rest, &models.InterpreterError{
				Message:        "expected binding clause after comma",
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]
	}

	if rest[0].Type != tokens.IN {
		return nil, rest, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: rest[0].SourceLocation,
		}
	}

	if len(rest) == 1 {
		return nil, rest, &models.InterpreterError{
			Message: "expected expression after \"in\"",
		}
	}
	rest = rest[1:]

	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &LetExpression{
		LetClauses: bindingExpressions,
		loc:        toks[0].SourceLocation,
		InClause:   exp2,
	}, rest, nil
}
