package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type FunctionExpression struct {
	Args    []types.Arg
	RetType types.Type
	body    expressions.Expression
	loc     *models.SourceLocation
}

type FuncValue struct {
	Bindings expressions.Bindings
	Exp      FunctionExpression
}

func (f *FuncValue) Call(args []any) (any, error) {
	if len(args) != len(f.Exp.Args) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected %d arguments, got %d", len(f.Exp.Args), len(args)),
			SourceLocation: f.Exp.loc,
		}
	}
	newBindings := make(expressions.Bindings)
	for k, v := range f.Bindings {
		newBindings[k] = v
	}
	for i, arg := range f.Exp.Args {
		newBindings[arg.Name] = args[i]
	}
	ret, err := f.Exp.body.Evaluate(newBindings)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (f *FuncValue) Args() []types.Arg {
	return f.Exp.Args
}

func (f *FuncValue) Return() types.Type {
	return f.Exp.RetType
}

func (f *FuncValue) String() string {
	return fmt.Sprintf("func(%v) %v { ... }", f.Exp.Args, f.Exp.RetType)
}

func (fe *FunctionExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	innerTB := make(types.TypeBindings)
	for k, v := range tb {
		innerTB[k] = v
	}

	argTypes := make([]types.Type, 0)
	for _, arg := range fe.Args {
		argTypes = append(argTypes, arg.Type)
		innerTB[arg.Name] = arg.Type
	}

	retType, err := fe.body.Type(innerTB)
	if err != nil {
		return nil, err
	}

	retSuper, innerErr := types.IsSuperTo(fe.RetType, retType)
	if innerErr != nil {
		return nil, &models.InterpreterError{
			Message:        "inconsistent return type",
			SourceLocation: fe.loc,
			Underlying:     innerErr,
		}
	}

	if !retSuper {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected return type %s, got %s", fe.RetType, retType),
			SourceLocation: fe.loc,
		}
	}

	return types.Func(argTypes, retType), nil
}

func (fe *FunctionExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
	// capture the current bindings
	retBindings := make(expressions.Bindings)
	for k, v := range bindings {
		retBindings[k] = v
	}

	return &FuncValue{
		Exp:      *fe,
		Bindings: retBindings,
	}, nil
}

func (fe *FunctionExpression) SourceLocation() *models.SourceLocation {
	return fe.loc
}

func parseFunction(toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
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
			SourceLocation: &tok.SourceLocation,
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
			SourceLocation: &tok.SourceLocation,
		}
	}
	toks.Pop()

	args := make([]types.Arg, 0)
	var popErr error
	for tok, popErr = toks.Pop(); popErr == nil; tok, popErr = toks.Pop() {
		if tok.Type == tokens.RIGHT_PAREN {
			break
		}

		if tok.Type != tokens.IDENTIFIER {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected identifier",
				SourceLocation: &tok.SourceLocation,
			}
		}
		argLoc := tok.SourceLocation
		argName := tok.Value

		tok, ok := toks.Peek()
		if !ok {
			return nil, &models.InterpreterError{
				Message:        "after argument declaration",
				SourceLocation: &argLoc,
				Underlying: &models.InterpreterError{
					Message:        "expected comma or closing parenthesis",
					Underlying:     popErr,
					SourceLocation: &tok.SourceLocation,
				},
			}
		}

		var argType types.Type = types.PrimitiveTypeAny
		if tok.Type != tokens.COMMA && tok.Type != tokens.RIGHT_PAREN {
			var innerErr error
			argType, innerErr = parseType(toks)
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "after argument declaration",
					SourceLocation: &argLoc,
					Underlying:     innerErr,
				}
			}
		}

		args = append(args, types.Arg{argName, argType})

		tok, innerErr := toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected comma or closing parenthesis",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			}
		}

		if tok.Type == tokens.RIGHT_PAREN {
			break
		}

		if tok.Type != tokens.COMMA {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected comma or closing parenthesis",
				SourceLocation: &tok.SourceLocation,
			}
		}
	}

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected return type",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	retType, innerErr := parseType(toks)
	if err != nil {
		return nil, &models.InterpreterError{
			Message:        "expected return type",
			SourceLocation: toks.CurrentSourceLocation(),
			Underlying:     innerErr,
		}
	}

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected function body",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	exp, err = ParseExpression(toks)
	if err != nil {
		return nil, err
	}

	return &FunctionExpression{
		Args:    args,
		RetType: retType,
		body:    exp,
		loc:     beginLoc,
	}, nil
}
