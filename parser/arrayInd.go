package parser

import (
	"errors"

	"github.com/brandonksides/phonk/models"
	"github.com/brandonksides/phonk/tokens"
)

type ArrayAccessExpression struct {
	Array models.Expression
	Index models.Expression
	loc   models.SourceLocation
}

func (aae *ArrayAccessExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	arr, err := aae.Array.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	arrSlice, ok := arr.([]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("expected array"),
			SourceLocation: aae.Array.SourceLocation(),
		}
	}

	index, err := aae.Index.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	indexInt, ok := index.(int)
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("expected int"),
			SourceLocation: aae.Index.SourceLocation(),
		}
	}

	if indexInt < 0 || indexInt >= len(arrSlice) {
		return nil, &models.InterpreterError{
			Err:            errors.New("index out of bounds"),
			SourceLocation: aae.Index.SourceLocation(),
		}
	}

	return arrSlice[indexInt], nil
}

func (aae *ArrayAccessExpression) SourceLocation() models.SourceLocation {
	return aae.loc
}

type ArraySliceExpression struct {
	Array models.Expression
	Begin *models.Expression
	End   *models.Expression
	loc   models.SourceLocation
}

func (ase *ArraySliceExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	arr, err := ase.Array.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	arrSlice, ok := arr.([]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Err:            errors.New("expected array"),
			SourceLocation: ase.Array.SourceLocation(),
		}
	}

	var beginInt int
	if ase.Begin != nil {
		begin, err := (*ase.Begin).Evaluate(bindings)
		if err != nil {
			return nil, err
		}
		beginInt, ok = begin.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: (*ase.Begin).SourceLocation(),
			}
		}
	}

	var endInt int = len(arrSlice)
	if ase.End != nil {
		end, err := (*ase.End).Evaluate(bindings)
		if err != nil {
			return nil, err
		}
		endInt, ok = end.(int)
		if !ok {
			return nil, &models.InterpreterError{
				Err:            errors.New("expected int"),
				SourceLocation: (*ase.End).SourceLocation(),
			}
		}
	}

	if beginInt < 0 || beginInt >= len(arrSlice) {
		return nil, &models.InterpreterError{
			Err:            errors.New("begin index out of bounds"),
			SourceLocation: (*ase.Begin).SourceLocation(),
		}
	}

	if endInt < 0 || endInt > len(arrSlice) {
		return nil, &models.InterpreterError{
			Err:            errors.New("end index out of bounds"),
			SourceLocation: (*ase.End).SourceLocation(),
		}
	}

	if beginInt > endInt {
		return nil, &models.InterpreterError{
			Err:            errors.New("begin index greater than end index"),
			SourceLocation: ase.SourceLocation(),
		}
	}

	return arrSlice[beginInt:endInt], nil
}

func (ase *ArraySliceExpression) SourceLocation() models.SourceLocation {
	return ase.loc
}

func parseArrayIndex(arr models.Expression, toks []tokens.Token) (exp models.Expression, rest []tokens.Token, err *models.InterpreterError) {
	if len(toks) == 0 {
		return nil, toks, &models.InterpreterError{
			Err: errors.New("expected token"),
		}
	}

	var idx *models.Expression
	rest = toks
	if toks[0].Type != tokens.COLON {
		var idxVal models.Expression
		idxVal, rest, err = ParseExpression(toks)
		if err != nil {
			return nil, toks, err
		}
		idx = &idxVal
	}

	if len(rest) == 0 {
		return nil, rest, &models.InterpreterError{
			Err: errors.New("unexpected end of input"),
		}
	}

	if rest[0].Type == tokens.COLON {
		rest = rest[1:]

		if len(rest) == 0 {
			return nil, rest, &models.InterpreterError{
				Err: errors.New("unexpected end of input"),
			}
		}

		var idx2 *models.Expression
		if rest[0].Type != tokens.RIGHT_SQUARE_BRACKET {
			var idxVal models.Expression
			idxVal, rest, err = ParseExpression(rest)
			if err != nil {
				return nil, rest, err
			}
			idx2 = &idxVal
		}

		return &ArraySliceExpression{
			Array: arr,
			Begin: idx,
			End:   idx2,
			loc:   arr.SourceLocation(),
		}, rest, nil
	}
	return &ArrayAccessExpression{
		Array: arr,
		Index: *idx,
		loc:   arr.SourceLocation(),
	}, rest, nil
}
