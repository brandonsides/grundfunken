package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type ObjectLiteralExpression struct {
	Fields map[string]expressions.Expression
	loc    *models.SourceLocation
}

func (ole *ObjectLiteralExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	fieldTypes := make(map[string]types.Type)
	for key, value := range ole.Fields {
		t, err := value.Type(tb)
		if err != nil {
			return nil, err
		}
		fieldTypes[key] = t
	}
	return types.Object(fieldTypes), nil
}

func (ole *ObjectLiteralExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	newBindings := make(map[string]any)
	for key, value := range bindings {
		newBindings[key] = value
	}

	obj := make(map[string]any)
	newBindings["this"] = obj

	for key, value := range ole.Fields {
		val, err := value.Evaluate(newBindings)
		if err != nil {
			return nil, err
		}
		obj[key] = val
	}
	return obj, nil
}

func (ole *ObjectLiteralExpression) SourceLocation() *models.SourceLocation {
	return ole.loc
}

func parseObjectLiteralExpression(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, innerErr := toks.Pop()
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "expected object literal expression",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_SQUIGGLY_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected object literal expression",
			SourceLocation: &tok.SourceLocation,
		}
	}

	fields := make(map[string]expressions.Expression)
	for {
		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message: "to terminate object literal expression",
				Underlying: &models.InterpreterError{
					Message:        "expected closing bracket",
					SourceLocation: toks.CurrentSourceLocation(),
					Underlying:     innerErr,
				},
			}
		}

		if tok.Type == tokens.RIGHT_SQUIGGLY_BRACKET {
			return &ObjectLiteralExpression{
				Fields: fields,
				loc:    beginLoc,
			}, nil
		}

		key := tok.Value
		keyLoc := tok.SourceLocation

		tok, innerErr = toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "in object literal",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "to bind object field " + key,
					SourceLocation: &keyLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected colon",
						SourceLocation: toks.CurrentSourceLocation(),
						Underlying:     innerErr,
					},
				},
			}
		}

		if tok.Type != tokens.COLON {
			return nil, &models.InterpreterError{
				Message:        "in object literal",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "to bind object field " + key,
					SourceLocation: &keyLoc,
					Underlying: &models.InterpreterError{
						Message:        "unexpected token; expected colon",
						SourceLocation: &tok.SourceLocation,
						Underlying:     innerErr,
					},
				},
			}
		}
		colLoc := tok.SourceLocation

		_, ok := toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "in object literal expression",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "to bind object field " + key,
					SourceLocation: &keyLoc,
					Underlying: &models.InterpreterError{
						Message:        "after colon",
						SourceLocation: &colLoc,
						Underlying: &models.InterpreterError{
							Message:        "expected expression",
							SourceLocation: toks.CurrentSourceLocation(),
							Underlying:     innerErr,
						},
					},
				},
			}
		}

		var exp1 expressions.Expression
		exp1, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}

		fields[key] = exp1

		tok, ok = toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "to terminate object literal expression",
				SourceLocation: beginLoc,
				Underlying: &models.InterpreterError{
					Message:        "expected closing bracket",
					SourceLocation: toks.CurrentSourceLocation(),
				},
			}
		}

		if tok.Type != tokens.COMMA {
			break
		}

		toks.Pop()
	}

	toks.Pop()

	if tok.Type != tokens.RIGHT_SQUIGGLY_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "to terminate object literal expression",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "unexpected token; expected closing bracket",
				SourceLocation: &tok.SourceLocation,
			},
		}
	}

	return &ObjectLiteralExpression{
		Fields: fields,
		loc:    beginLoc,
	}, nil
}
