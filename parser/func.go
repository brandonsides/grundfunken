package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

type FunctionExpression struct {
	args []string
	exp  models.Expression
}

func (fe *FunctionExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	retBindings := make(models.Bindings)
	for k, v := range bindings {
		retBindings[k] = v
	}

	return &models.ExpFunction{
		Args:     fe.args,
		Bindings: retBindings,
		Exp:      fe.exp,
	}, nil
}

func (fe *FunctionExpression) SourceLocation() models.SourceLocation {
	return fe.exp.SourceLocation()
}

func parseFunction(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.FUNC {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.LEFT_PAREN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	args := make([]string, 0)
	for len(rest) > 0 {
		if rest[0].Type == tokens.RIGHT_PAREN {
			break
		}

		if rest[0].Type != tokens.IDENTIFIER {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		args = append(args, rest[0].Value)
		rest = rest[1:]
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		if rest[0].Type == tokens.RIGHT_PAREN {
			break
		}

		if rest[0].Type != tokens.COMMA {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		rest = rest[1:]
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	rest = rest[1:]

	exp, rest, err = ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &FunctionExpression{
		args: args,
		exp:  exp,
	}, rest, nil
}
