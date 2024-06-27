package parser

import (
	"strconv"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseAtomic(toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	tok := toks.Peek()
	if tok == nil {
		return nil, &models.InterpreterError{
			Message:        "expected token",
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
		tok = toks.Pop()
		if tok.Type != tokens.RIGHT_PAREN {
			return nil, &models.InterpreterError{
				Message:        "expected closing parenthesis",
				SourceLocation: tok.SourceLocation,
			}
		}
	case tokens.NUMBER, tokens.MINUS:
		var numStr string

		tok = toks.Peek()
		if tok.Type == tokens.MINUS {
			toks.Pop()

			numStr = "-"
		}

		tok = toks.Pop()
		if tok == nil {
			return nil, &models.InterpreterError{
				Message:        "expected token",
				SourceLocation: toks.CurrentSourceLocation(),
			}
		}
		if tok.Type != tokens.NUMBER {
			return nil, &models.InterpreterError{
				Message:        "unexpected token",
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
		if tok.Value == "true" {
			exp, err = &LiteralExpression{
				val: true,
				loc: toks[0].SourceLocation,
			}, nil
		} else if toks[0].Value == "false" {
			exp, rest, err = &LiteralExpression{
				val: false,
				loc: toks[0].SourceLocation,
			}, toks[1:], nil
		} else {
			exp, rest, err = &IdentifierExpression{
				name: toks[0].Value,
				loc:  toks[0].SourceLocation,
			}, toks[1:], nil
		}
	case tokens.LEFT_SQUARE_BRACKET:
		exp, rest, err = parseArrayLiteral(toks)
	case tokens.LEFT_SQUIGGLY_BRACKET:
		exp, rest, err = parseObjectLiteralExpression(toks)
	case tokens.STRING:
		exp, rest, err = &LiteralExpression{
			val: toks[0].Value,
			loc: toks[0].SourceLocation,
		}, toks[1:], nil
	case tokens.LET:
		exp, rest, err = parseLetExpression(toks)
	case tokens.IF:
		exp, rest, err = parseIfExpression(toks)
	default:
		return nil, toks, nil
	}
	if err != nil {
		return nil, rest, err
	}

	for len(rest) != 0 && (rest[0].Type == tokens.LEFT_SQUARE_BRACKET || rest[0].Type == tokens.LEFT_PAREN || rest[0].Type == tokens.DOT) {
		tok := rest[0]
		rest = rest[1:]

		switch tok.Type {
		case tokens.LEFT_PAREN:
			var exps []models.Expression
			exps, rest, err = parseExpressions(rest)
			if err != nil {
				return nil, rest, err
			}
			if rest[0].Type != tokens.RIGHT_PAREN {
				return nil, rest, &models.InterpreterError{
					Message:        "unexpected token",
					SourceLocation: rest[0].SourceLocation,
				}
			}
			rest = rest[1:]
			exp = &FunctionCallExpression{
				Function: exp,
				Args:     exps,
				loc:      toks[0].SourceLocation,
			}
		case tokens.LEFT_SQUARE_BRACKET:
			exp, rest, err = parseArrayIndex(exp, rest)
			if err != nil {
				return nil, rest, err
			}
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Message: "unexpected end of input",
				}
			}
			if rest[0].Type != tokens.RIGHT_SQUARE_BRACKET {
				return nil, rest, &models.InterpreterError{
					Message:        "unexpected token",
					SourceLocation: rest[0].SourceLocation,
				}
			}
			rest = rest[1:]
		case tokens.DOT:
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Message:        "unexpected end of input",
					SourceLocation: tok.SourceLocation,
				}
			}
			if rest[0].Type != tokens.IDENTIFIER {
				return nil, rest, &models.InterpreterError{
					Message:        "unexpected token",
					SourceLocation: rest[0].SourceLocation,
				}
			}
			exp = &FieldAccessExpression{
				Object: exp,
				Field:  rest[0].Value,
				loc: models.SourceLocation{
					LineNumber:   exp.SourceLocation().LineNumber,
					ColumnNumber: exp.SourceLocation().ColumnNumber,
				},
			}
			rest = rest[1:]
		}
	}

	if len(rest) != 0 && rest[0].Type == tokens.FOR {
		exp, rest, err = parseForExpression(exp, rest)
	}

	return exp, rest, err
}
