package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type FunctionExpression struct {
	args []string
	body models.Expression
	loc  models.SourceLocation
}

type FuncValue struct {
	Bindings models.Bindings
	Exp      FunctionExpression
}

func (f FuncValue) Call(args []any) (any, error) {
	if len(args) != len(f.Exp.args) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected %d arguments, got %d", len(f.Exp.args), len(args)),
			SourceLocation: f.Exp.loc,
		}
	}
	newBindings := make(models.Bindings)
	for k, v := range f.Bindings {
		newBindings[k] = v
	}
	for i, arg := range f.Exp.args {
		newBindings[arg] = args[i]
	}
	ret, err := f.Exp.body.Evaluate(newBindings)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (fe *FunctionExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	// capture the current bindings
	retBindings := make(models.Bindings)
	for k, v := range bindings {
		retBindings[k] = v
	}

	return &FuncValue{
		Exp:      *fe,
		Bindings: retBindings,
	}, nil
}

func (fe *FunctionExpression) SourceLocation() models.SourceLocation {
	return fe.loc
}

func parseFunction(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	if toks[0].Type != tokens.FUNC {
		return nil, toks, &models.InterpreterError{
			Message:        "unexpected token",
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	if rest[0].Type != tokens.LEFT_PAREN {
		return nil, rest, &models.InterpreterError{
			Message:        "unexpected token",
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
				Message:        "unexpected token",
				SourceLocation: rest[0].SourceLocation,
			}
		}

		args = append(args, rest[0].Value)
		rest = rest[1:]
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Message: "unexpected end of input",
			}
		}

		if rest[0].Type == tokens.RIGHT_PAREN {
			break
		}

		if rest[0].Type != tokens.COMMA {
			return nil, rest, &models.InterpreterError{
				Message:        "unexpected token",
				SourceLocation: rest[0].SourceLocation,
			}
		}

		rest = rest[1:]
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Message: "unexpected end of input",
		}
	}

	rest = rest[1:]

	exp, rest, err = ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &FunctionExpression{
		args: args,
		body: exp,
		loc:  toks[0].SourceLocation,
	}, rest, nil
}
