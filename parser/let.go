package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type BindingExpression struct {
	Identifier      string
	ExpectedType    types.Type
	ExpectedTypeLoc *models.SourceLocation
	Expression      expressions.Expression
}

type LetExpression struct {
	loc        *models.SourceLocation
	LetClauses []BindingExpression
	InClause   expressions.Expression
}

func (le *LetExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	newTB := make(types.TypeBindings)
	for k, v := range tb {
		newTB[k] = v
	}

	for _, bindingExp := range le.LetClauses {
		if funcExp, ok := bindingExp.Expression.(*FunctionExpression); ok {
			typs := make([]types.Type, 0, len(funcExp.Args))
			for _, arg := range funcExp.Args {
				typs = append(typs, arg.Type)
			}
			newTB[bindingExp.Identifier] = types.Func(typs, funcExp.RetType)
		}

		t, err := bindingExp.Expression.Type(newTB)
		if err != nil {
			return nil, err
		}

		if bindingExp.ExpectedTypeLoc != nil {
			isSuper, innerErr := types.IsSuperTo(bindingExp.ExpectedType, t)
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "in let clause",
					SourceLocation: le.SourceLocation(),
					Underlying: &models.InterpreterError{
						Message:        "error checking type constraint",
						SourceLocation: bindingExp.ExpectedTypeLoc,
						Underlying: &models.InterpreterError{
							Message:        "on expression",
							SourceLocation: bindingExp.Expression.SourceLocation(),
							Underlying:     innerErr,
						},
					},
				}
			}

			if !isSuper {
				return nil, &models.InterpreterError{
					Message:        "in let clause",
					SourceLocation: le.SourceLocation(),
					Underlying: &models.InterpreterError{
						Message:        "unmet type constraint",
						SourceLocation: bindingExp.ExpectedTypeLoc,
						Underlying: &models.InterpreterError{
							Message:        "on expression",
							SourceLocation: bindingExp.Expression.SourceLocation(),
						},
					},
				}
			}

			newTB[bindingExp.Identifier] = bindingExp.ExpectedType
		} else {
			newTB[bindingExp.Identifier] = t
		}
	}

	return le.InClause.Type(newTB)
}

func (le *LetExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	newBindings := make(expressions.Bindings)
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

func (le *LetExpression) SourceLocation() *models.SourceLocation {
	return le.loc
}

func parseLetExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected let expression",
			SourceLocation: beginLoc,
			Underlying:     innerErr,
		}
	}

	if tok.Type != tokens.LET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected let clause",
			SourceLocation: &tok.SourceLocation,
		}
	}

	tok, innerErr = toks.Pop()
	if innerErr != nil {
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
					SourceLocation: &tok.SourceLocation,
				},
				SourceLocation: beginLoc,
			}
		}
		identifier := tok.Value
		identifierDeclLoc := tok.SourceLocation

		var ok bool
		tok, ok = toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message: "in let clause",
				Underlying: &models.InterpreterError{
					Message:        "expected equal sign or type constraint",
					SourceLocation: toks.CurrentSourceLocation(),
					Underlying:     innerErr,
				},
				SourceLocation: beginLoc,
			}
		}

		var typ types.Type = types.PrimitiveTypeAny
		var typLoc *models.SourceLocation
		if tok.Type != tokens.EQUAL {
			loc := tok.SourceLocation
			typLoc = &loc
			typ, innerErr = parseType(toks)
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message: "in let clause",
					Underlying: &models.InterpreterError{
						Message:        fmt.Sprintf("unexpected token \"%s\"; expected equal sign or type constraint", tok.Value),
						SourceLocation: &tok.SourceLocation,
					},
					SourceLocation: beginLoc,
				}
			}
		}

		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "in let clause",
				SourceLocation: beginLoc,
				Underlying:     innerErr,
			}
		}
		if tok.Type != tokens.EQUAL {
			return nil, &models.InterpreterError{
				Message:        "in let clause",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "expected equal sign",
					SourceLocation: &tok.SourceLocation,
				},
			}
		}

		tok, ok = toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message: "in let clause",
				Underlying: &models.InterpreterError{
					Message:        "in binding clause for identifier \"" + identifier + "\"",
					SourceLocation: &identifierDeclLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected expression",
						SourceLocation: toks.CurrentSourceLocation(),
					},
				},
				SourceLocation: beginLoc,
			}
		}

		var exp1 expressions.Expression
		exp1, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}

		bindingExpressions = append(bindingExpressions, BindingExpression{
			Identifier:      identifier,
			Expression:      exp1,
			ExpectedType:    typ,
			ExpectedTypeLoc: typLoc,
		})

		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "in let clause",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "after binding clause for identifier \"" + identifier + "\"",
					SourceLocation: &identifierDeclLoc,
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
					SourceLocation: &tok.SourceLocation,
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
				SourceLocation: &tok.SourceLocation,
			},
		}
	}

	_, ok := toks.Peek()
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
