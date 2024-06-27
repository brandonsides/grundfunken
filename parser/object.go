package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type ObjectLiteralExpression struct {
	Fields map[string]models.Expression
	loc    models.SourceLocation
}

func (ole *ObjectLiteralExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	obj := make(map[string]any)
	for key, value := range ole.Fields {
		val, err := value.Evaluate(bindings)
		if err != nil {
			return nil, err
		}
		obj[key] = val
	}
	return obj, nil
}

func (ole *ObjectLiteralExpression) SourceLocation() models.SourceLocation {
	return ole.loc
}

func parseObjectLiteralExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) < 2 {
		return nil, toks, &models.InterpreterError{
			Message:        "unexpected end of input",
			SourceLocation: toks[0].SourceLocation,
		}
	}

	if toks[0].Type != tokens.LEFT_SQUIGGLY_BRACKET {
		return nil, toks, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	fields := make(map[string]models.Expression)
	for {
		if len(rest) == 0 {
			return nil, toks, &models.InterpreterError{
				Message:        "unexpected end of input",
				SourceLocation: toks[0].SourceLocation,
			}
		}

		if rest[0].Type == tokens.RIGHT_SQUIGGLY_BRACKET {
			return &ObjectLiteralExpression{
				Fields: fields,
				loc:    toks[0].SourceLocation,
			}, rest[1:], nil
		}

		key := rest[0].Value
		rest = rest[1:]

		if len(rest) == 0 {
			return nil, toks, &models.InterpreterError{
				Message:        "unexpected end of input",
				SourceLocation: toks[0].SourceLocation,
			}
		}

		if rest[0].Type != tokens.COLON {
			return nil, toks, &models.InterpreterError{
				Message:        "unexpected token",
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]

		if len(rest) == 0 {
			return nil, toks, &models.InterpreterError{
				Message: "unexpected end of input",
			}
		}

		var exp1 models.Expression
		exp1, rest, err = ParseExpression(rest)
		if err != nil {
			return nil, toks, &models.InterpreterError{
				Message:        "in object literal expression",
				Underlying:     err,
				SourceLocation: rest[0].SourceLocation,
			}
		}

		fields[key] = exp1

		if len(rest) == 0 {
			return nil, toks, &models.InterpreterError{
				Message: "unexpected end of input",
			}
		}

		if rest[0].Type != tokens.COMMA {
			break
		}

		rest = rest[1:]
	}

	if len(rest) == 0 {
		return nil, toks, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	if rest[0].Type != tokens.RIGHT_SQUIGGLY_BRACKET {
		return nil, toks, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: rest[0].SourceLocation,
		}
	}

	return &ObjectLiteralExpression{
		Fields: fields,
		loc:    toks[0].SourceLocation,
	}, rest[1:], nil
}
