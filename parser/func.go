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

func parseFunction(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected function declaration",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.FUNC {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected function declaration",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected argument list",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_PAREN {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected left parenthesis",
			SourceLocation: tok.SourceLocation,
		}
	}
	toks.Pop()

	args := make([]string, 0)
	var popErr error
	for tok, popErr = toks.Pop(); popErr == nil; tok, popErr = toks.Pop() {
		if tok.Type == tokens.RIGHT_PAREN {
			break
		}

		if tok.Type != tokens.IDENTIFIER {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected identifier",
				SourceLocation: tok.SourceLocation,
			}
		}
		args = append(args, tok.Value)

		tok, popErr := toks.Pop()
		if popErr != nil {
			return nil, &models.InterpreterError{
				Message:        "after comma",
				SourceLocation: tok.SourceLocation,
				Underlying: &models.InterpreterError{
					Message:        "expected new argument",
					Underlying:     popErr,
					SourceLocation: tok.SourceLocation,
				},
			}
		}

		if tok.Type == tokens.RIGHT_PAREN {
			break
		}

		if tok.Type != tokens.COMMA {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected comma or closing parenthesis",
				SourceLocation: tok.SourceLocation,
			}
		}
	}

	exp, err = ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	return &FunctionExpression{
		args: args,
		body: exp,
		loc:  beginLoc,
	}, nil
}
