package parser

import (
	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseType(toks *tokens.TokenStack) (types.Type, error) {
	return parseSumType(toks)
}

func parseSumType(toks *tokens.TokenStack) (types.Type, error) {
	t1, err := parseFuncType(toks)
	if err != nil {
		return nil, err
	}

	tok, ok := toks.Peek()
	if !ok {
		return t1, nil
	}

	if tok.Type != tokens.PIPE {
		return t1, nil
	}
	toks.Pop()

	t2, err := parseSumType(toks)
	if err != nil {
		return nil, err
	}

	return types.Sum(t1, t2), nil
}

func parseFuncType(toks *tokens.TokenStack) (types.Type, error) {
	tok, ok := toks.Peek()
	if !ok {
		return types.PrimitiveTypeAny, nil
	}

	if tok.Type != tokens.FUNC {
		return parseAtomicType(toks)
	}
	toks.Pop()

	tok, err := toks.Pop()
	if err != nil {
		return nil, &models.InterpreterError{
			Message:        "expected opening parenthesis",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_PAREN {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected opening parenthesis",
			SourceLocation: &tok.SourceLocation,
		}
	}

	args := make([]types.Type, 0)
	for {
		argType, err := parseType(toks)
		if err != nil {
			return nil, err
		}

		args = append(args, argType)

		tok, err = toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected comma or closing parenthesis",
				SourceLocation: toks.CurrentSourceLocation(),
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

	retType, err := parseFuncType(toks)
	if err != nil {
		return nil, err
	}

	return types.Func(args, retType), nil
}

func parseAtomicType(toks *tokens.TokenStack) (types.Type, error) {
	tok, ok := toks.Peek()
	if !ok {
		return types.PrimitiveTypeAny, nil
	}

	switch tok.Type {
	case tokens.FUNC:
		return parseFuncType(toks)
	case tokens.LEFT_PAREN:
		toks.Pop()
		typ, err := parseType(toks)
		if err != nil {
			return nil, err
		}
		tok, err := toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected closing parenthesis",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}

		if tok.Type != tokens.RIGHT_PAREN {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected closing parenthesis",
				SourceLocation: &tok.SourceLocation,
			}
		}

		return typ, nil
	case tokens.IDENTIFIER:
		tok, innerErr := toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected identifier",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			}
		}

		return types.ParsePrimitive(tok.Value), nil
	case tokens.LEFT_SQUARE_BRACKET:
		toks.Pop()
		typ, err := parseType(toks)
		if err != nil {
			return nil, err
		}
		tok, err := toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected closing square bracket",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}

		if tok.Type != tokens.RIGHT_SQUARE_BRACKET {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected closing square bracket",
				SourceLocation: &tok.SourceLocation,
			}
		}

		return types.List(typ), nil
	case tokens.LEFT_SQUIGGLY_BRACKET:
		typ, err := parseObjectType(toks)
		if err != nil {
			return nil, err
		}

		return typ, nil
	default:
		return types.PrimitiveTypeAny, nil
	}
}

func parseObjectType(toks *tokens.TokenStack) (types.Type, error) {
	tok, err := toks.Pop()
	if err != nil {
		return nil, &models.InterpreterError{
			Message:        "expected opening squiggly bracket",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	if tok.Type != tokens.LEFT_SQUIGGLY_BRACKET {
		return nil, &models.InterpreterError{
			Message:        "unexpected token; expected opening squiggly bracket",
			SourceLocation: &tok.SourceLocation,
		}
	}

	fields := make(map[string]types.Type)
	for {
		tok, err = toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected identifier or closing squiggly bracket",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}

		if tok.Type == tokens.RIGHT_SQUIGGLY_BRACKET {
			break
		}

		if tok.Type != tokens.IDENTIFIER {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected identifier or closing squiggly bracket",
				SourceLocation: &tok.SourceLocation,
			}
		}
		fieldName := tok.Value

		tok, err = toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected colon",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}

		if tok.Type != tokens.COLON {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected colon",
				SourceLocation: &tok.SourceLocation,
			}
		}

		typ, err := parseType(toks)
		if err != nil {
			return nil, err
		}

		fields[fieldName] = typ

		tok, err = toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected comma or closing squiggly bracket",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}

		if tok.Type == tokens.RIGHT_SQUIGGLY_BRACKET {
			break
		}

		if tok.Type != tokens.COMMA {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected comma or closing squiggly bracket",
				SourceLocation: &tok.SourceLocation,
			}
		}
	}

	return types.Object(fields), nil
}
