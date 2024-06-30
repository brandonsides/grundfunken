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

func parseLetExpression(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected let expression",
			SourceLocation: beginLoc,
		}
	}

	if tok.Type != tokens.LET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected let clause",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	tok, innerErr := toks.Pop()
	if err != nil {
		return nil, &models.InterpreterError{
			Message: "in let clause",
			Underlying: &models.InterpreterError{
				Message:        "expected identifier",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			},
			SourceLocation: beginLoc,
		}
	}

	bindingExpressions := make([]BindingExpression, 0)
	for {
		if tok.Type != tokens.IDENTIFIER {
			return nil, &models.InterpreterError{
				Message: "in let clause",
				Underlying: &models.InterpreterError{
					Message:        "unexpected token; expected identifier",
					SourceLocation: tok.SourceLocation,
				},
				SourceLocation: beginLoc,
			}
		}
		identifier := tok.Value
		identifierDeclLoc := tok.SourceLocation

		tok, innerErr = toks.Pop()
		if tok.Type != tokens.EQUAL {
			return nil, &models.InterpreterError{
				Message: "in let clause",
				Underlying: &models.InterpreterError{
					Message:        "expected equal sign",
					SourceLocation: tok.SourceLocation,
					Underlying:     innerErr,
				},
			}
		}

		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message: "in let clause",
				Underlying: &models.InterpreterError{
					Message:        "in binding clause for identifier \"" + identifier + "\"",
					SourceLocation: identifierDeclLoc,
					Underlying: &models.InterpreterError{
						Message:        fmt.Sprintf("expected expression"),
						SourceLocation: toks.CurrentSourceLocation(),
						Underlying:     innerErr,
					},
				},
				SourceLocation: beginLoc,
			}
		}

		var exp1 models.Expression
		exp1, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}

		bindingExpressions = append(bindingExpressions, BindingExpression{
			Identifier: identifier,
			Expression: exp1,
		})

		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "in let clause",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "after binding clause for identifier \"" + identifier + "\"",
					SourceLocation: identifierDeclLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected \"in\" clause",
						Underlying:     innerErr,
						SourceLocation: exp1.SourceLocation(),
					},
				},
			}
		}

		if tok.Type != tokens.COMMA {
			break
		}

		if tok, innerErr = toks.Pop(); innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "in let clause",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "after comma",
					SourceLocation: tok.SourceLocation,
					Underlying: &models.InterpreterError{
						Message:        "expected identifier for next binding",
						SourceLocation: toks.CurrentSourceLocation(),
						Underlying:     innerErr,
					},
				},
			}
		}
	}

	if tok.Type != tokens.IN {
		return nil, &models.InterpreterError{
			Message:        "in let clause",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected \"in\" clause",
				SourceLocation: tok.SourceLocation,
			},
		}
	}

	_, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "in let clause",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "expected expression after \"in\"",
				SourceLocation: toks.CurrentSourceLocation(),
			},
		}
	}

	exp2, err := ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	return &LetExpression{
		LetClauses: bindingExpressions,
		loc:        beginLoc,
		InClause:   exp2,
	}, nil
}
