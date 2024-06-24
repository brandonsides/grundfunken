package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

func parseAtomic(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	switch toks[0].Type {
	case tokens.FUNC:
		exp, rest, err = parseFunction(toks)
	case tokens.LEFT_PAREN:
		rest = toks[1:]
		exp, rest, err = ParseExpression(rest)
		if err != nil {
			return nil, rest, err
		}
		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("expected closing parenthesis"),
			}
		}
		if rest[0].Type != tokens.RIGHT_PAREN {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("expected closing parenthesis"),
				SourceLocation: rest[0].SourceLocation,
			}
		}
		rest = rest[1:]
	case tokens.NUMBER, tokens.MINUS:
		var numStr string
		rest = toks
		if toks[0].Type == tokens.MINUS {
			rest = toks[1:]
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Err: errors.New("unexpected end of input"),
				}
			}

			numStr = "-"
		}

		if rest[0].Type != tokens.NUMBER {
			return nil, rest, &models.InterpreterError{
				Err:            errors.New("unexpected token"),
				SourceLocation: rest[0].SourceLocation,
			}
		}

		numStr += rest[0].Value
		rest = rest[1:]

		ret, innerErr := strconv.Atoi(numStr)
		if innerErr != nil {
			return nil, toks[1:], &models.InterpreterError{
				Err:            fmt.Errorf("failed to parse number literal: %w", innerErr),
				SourceLocation: toks[0].SourceLocation,
			}
		}

		exp = &LiteralExpression{
			val: ret,
			loc: toks[0].SourceLocation,
		}
	case tokens.IDENTIFIER:
		if toks[0].Value == "true" {
			exp, rest, err = &LiteralExpression{
				val: true,
				loc: toks[0].SourceLocation,
			}, toks[1:], nil
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
					Err:            errors.New("unexpected token"),
					SourceLocation: rest[0].SourceLocation,
				}
			}
			rest = rest[1:]
			exp = &FunctionCallExpression{
				Function: exp,
				Args:     exps,
				loc: models.SourceLocation{
					LineNumber:   exp.SourceLocation().LineNumber,
					ColumnNumber: exp.SourceLocation().ColumnNumber,
				},
			}
		case tokens.LEFT_SQUARE_BRACKET:
			exp, rest, err = parseArrayIndex(exp, rest)
			if err != nil {
				return nil, rest, err
			}
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Err: errors.New("unexpected end of input"),
				}
			}
			if rest[0].Type != tokens.RIGHT_SQUARE_BRACKET {
				return nil, rest, &models.InterpreterError{
					Err:            errors.New("unexpected token"),
					SourceLocation: rest[0].SourceLocation,
				}
			}
			rest = rest[1:]
		case tokens.DOT:
			if len(rest) == 0 {
				return nil, rest, &models.InterpreterError{
					Err: errors.New("unexpected end of input"),
				}
			}
			if rest[0].Type != tokens.IDENTIFIER {
				return nil, rest, &models.InterpreterError{
					Err:            errors.New("unexpected token"),
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
