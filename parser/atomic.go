package parser

import (
	"strconv"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseAtomic(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected expression",
			SourceLocation: toks.CurrentSourceLocation(),
		}
	}

	switch tok.Type {
	case tokens.FUNC:
		exp, err = parseFunction(toks)
	case tokens.LEFT_PAREN:
		toks.Pop()
		exp, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}
		tok, err := toks.Pop()
		if err != nil {
			return nil, &models.InterpreterError{
				Message:        "expected closing parenthesis",
				SourceLocation: exp.SourceLocation(),
			}
		}

		if tok.Type != tokens.RIGHT_PAREN {
			return nil, &models.InterpreterError{
				Message:        "expected closing parenthesis",
				SourceLocation: tok.SourceLocation,
			}
		}
	case tokens.NUMBER, tokens.MINUS:
		var numStr string

		tok, innerErr := toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected number",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			}
		}

		if tok.Type == tokens.MINUS {
			numStr = "-"
			tok, innerErr = toks.Pop()
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "expected number",
					SourceLocation: toks.CurrentSourceLocation(),
					Underlying:     innerErr,
				}
			}
		}

		if tok.Type != tokens.NUMBER {
			return nil, &models.InterpreterError{
				Message:        "unexpected token; expected number",
				SourceLocation: tok.SourceLocation,
			}
		}

		numStr += tok.Value

		ret, innerErr := strconv.Atoi(numStr)
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "failed to parse number literal",
				Underlying:     innerErr,
				SourceLocation: tok.SourceLocation,
			}
		}

		exp = &LiteralExpression{
			val: ret,
			loc: tok.SourceLocation,
		}
	case tokens.IDENTIFIER:
		tok, innerErr := toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected identifier",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			}
		}

		if tok.Value == "true" {
			exp, err = &LiteralExpression{
				val: true,
				loc: tok.SourceLocation,
			}, nil
		} else if tok.Value == "false" {
			exp, err = &LiteralExpression{
				val: false,
				loc: tok.SourceLocation,
			}, nil
		} else {
			exp, err = &IdentifierExpression{
				name: tok.Value,
				loc:  tok.SourceLocation,
			}, nil
		}
	case tokens.LEFT_SQUARE_BRACKET:
		exp, err = parseArrayLiteral(toks)
	case tokens.LEFT_SQUIGGLY_BRACKET:
		exp, err = parseObjectLiteralExpression(toks)
	case tokens.STRING:
		tok, innerErr := toks.Pop()
		if innerErr != nil {
			return nil, &models.InterpreterError{
				Message:        "expected string",
				SourceLocation: toks.CurrentSourceLocation(),
				Underlying:     innerErr,
			}
		}
		exp, err = &LiteralExpression{
			val: tok.Value,
			loc: tok.SourceLocation,
		}, nil
	case tokens.LET:
		exp, err = parseLetExpression(toks)
	case tokens.IF:
		exp, err = parseIfExpression(toks)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	for tok, ok := toks.Peek(); ok && (tok.Type == tokens.LEFT_SQUARE_BRACKET || tok.Type == tokens.LEFT_PAREN || tok.Type == tokens.DOT); tok, ok = toks.Peek() {
		switch tok.Type {
		case tokens.LEFT_PAREN:
			parenLoc := tok.SourceLocation
			var exps []models.Expression

			toks.Pop()
			exps, err = parseExpressions(toks)
			if err != nil {
				return nil, err
			}

			tok, innerErr := toks.Pop()
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "to terminate function call",
					SourceLocation: parenLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected closing parenthesis",
						SourceLocation: toks.CurrentSourceLocation(),
						Underlying:     innerErr,
					},
				}
			}
			if tok.Type != tokens.RIGHT_PAREN {
				return nil, &models.InterpreterError{
					Message:        "to terminate function call",
					SourceLocation: parenLoc,
					Underlying: &models.InterpreterError{
						Message:        "unexpected token; expected closing parenthesis",
						SourceLocation: tok.SourceLocation,
					},
				}
			}
			exp = &FunctionCallExpression{
				Function: exp,
				Args:     exps,
				loc:      beginLoc,
			}
		case tokens.LEFT_SQUARE_BRACKET:
			bracketLoc := tok.SourceLocation

			toks.Pop()
			exp, err = parseArrayIndex(exp, toks)
			if err != nil {
				return nil, err
			}

			tok, innerErr := toks.Pop()
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "to terminate array index",
					SourceLocation: bracketLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected closing square bracket",
						SourceLocation: exp.SourceLocation(),
					},
				}
			}
			if tok.Type != tokens.RIGHT_SQUARE_BRACKET {
				return nil, &models.InterpreterError{
					Message:        "to terminate array index",
					SourceLocation: bracketLoc,
					Underlying: &models.InterpreterError{
						Message:        "unexpected token; expected closing square bracket",
						SourceLocation: tok.SourceLocation,
					},
				}
			}
		case tokens.DOT:
			dotLoc := tok.SourceLocation
			toks.Pop()

			tok, innerErr := toks.Pop()
			if innerErr != nil {
				return nil, &models.InterpreterError{
					Message:        "in object field access",
					SourceLocation: dotLoc,
					Underlying: &models.InterpreterError{
						Message:        "expected identifier",
						SourceLocation: toks.CurrentSourceLocation(),
						Underlying:     innerErr,
					},
				}
			}

			if tok.Type != tokens.IDENTIFIER {
				return nil, &models.InterpreterError{
					Message:        "unexpected token; expected identifier",
					SourceLocation: tok.SourceLocation,
				}
			}
			exp = &FieldAccessExpression{
				Object: exp,
				Field:  tok.Value,
				loc:    exp.SourceLocation(),
			}
		}
	}

	tok, ok = toks.Peek()
	if ok && tok.Type == tokens.FOR {
		exp, err = parseForExpression(exp, toks)
		if err != nil {
			return nil, err
		}
	}

	return exp, err
}
