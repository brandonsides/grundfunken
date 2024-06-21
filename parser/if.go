package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

type IfExpression struct {
	Condition models.Expression
	Then      models.Expression
	Else      models.Expression
	loc       models.SourceLocation
}

func (ie *IfExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	cond, err := ie.Condition.Evaluate(bindings)
	if err != nil {
		return nil, err
	}

	condBool, ok := cond.(bool)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("if condition must evaluate to a boolean"),
			SourceLocation: ie.Condition.SourceLocation(),
		}
	}

	if condBool {
		return ie.Then.Evaluate(bindings)
	}

	return ie.Else.Evaluate(bindings)
}

func (ie *IfExpression) SourceLocation() models.SourceLocation {
	return ie.loc
}

func parseIfExpression(toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if toks[0].Type != tokens.IF {
		return nil, toks, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: toks[0].SourceLocation,
		}
	}

	rest = toks[1:]
	exp1, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.THEN {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp2, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type != tokens.ELSE {
		return nil, rest, &models.InterpreterError{
			Err:            errors.New("unexpected token"),
			SourceLocation: rest[0].SourceLocation,
		}
	}

	rest = rest[1:]
	exp3, rest, err := ParseExpression(rest)
	if err != nil {
		return nil, rest, err
	}

	return &IfExpression{
		Condition: exp1,
		Then:      exp2,
		Else:      exp3,
		loc:       toks[0].SourceLocation,
	}, rest, nil
}
