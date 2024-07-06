package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/tokens"
)

type ArrayAccessExpression struct {
	Array models.Expression
	Index models.Expression
	loc   models.SourceLocation
}

func (aae *ArrayAccessExpression) Type() (models.Type, *models.InterpreterError) {
	return models.PrimitiveTypeAny, nil
}

func (aae *ArrayAccessExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	arr, err := aae.Array.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	arrSlice, ok := arr.([]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected array; got %v", arr),
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
			Message:        fmt.Sprintf("expected int; got %v", index),
			SourceLocation: aae.Index.SourceLocation(),
		}
	}

	if indexInt < 0 || indexInt >= len(arrSlice) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("index out of bounds (%d)", index),
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

func (ase *ArraySliceExpression) Type() (models.Type, *models.InterpreterError) {
	_, err := ase.Array.Type()
	if err != nil {
		return nil, err
	}

	if ase.Begin != nil {
		t, err := (*ase.Begin).Type()
		if err != nil {
			return nil, err
		}
		if t != models.PrimitiveTypeInt {
			return nil, &models.InterpreterError{
				Message:        fmt.Sprintf("expected int; got %s", t),
				SourceLocation: (*ase.Begin).SourceLocation(),
			}
		}
	}

	return models.PrimitiveTypeList, nil
}

func (ase *ArraySliceExpression) Evaluate(bindings models.Bindings) (any, *models.InterpreterError) {
	arr, err := ase.Array.Evaluate(bindings)
	if err != nil {
		return nil, err
	}
	arrSlice, ok := arr.([]interface{})
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected array; got %v", arr),
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
				Message:        fmt.Sprintf("expected int; got %v", begin),
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
				Message:        fmt.Sprintf("expected int; got %v", end),
				SourceLocation: (*ase.End).SourceLocation(),
			}
		}
	}

	if beginInt < 0 || beginInt > len(arrSlice) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("begin index out of bounds (%d)", beginInt),
			SourceLocation: (*ase.Begin).SourceLocation(),
		}
	}

	if endInt < 0 || endInt > len(arrSlice) {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("end index out of bounds (%d)", endInt),
			SourceLocation: (*ase.End).SourceLocation(),
		}
	}

	if beginInt > endInt {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("begin index %d greater than end index %d", beginInt, endInt),
			SourceLocation: ase.SourceLocation(),
		}
	}

	return arrSlice[beginInt:endInt], nil
}

func (ase *ArraySliceExpression) SourceLocation() models.SourceLocation {
	return ase.loc
}

func parseArrayIndex(arr models.Expression, toks *tokens.TokenStack) (exp models.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected array index expression",
			SourceLocation: arr.SourceLocation(),
		}
	}

	var idx *models.Expression
	if tok.Type != tokens.COLON {
		var idxVal models.Expression
		idxVal, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}
		if idxVal != nil {
			idx = &idxVal
		}
	}

	tok, ok = toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "to terminate array index",
			SourceLocation: beginLoc,
			Underlying: &models.InterpreterError{
				Message:        "expected right square bracket",
				SourceLocation: arr.SourceLocation(),
			},
		}
	}

	if tok.Type == tokens.COLON {
		toks.Pop()

		var idx2 *models.Expression
		var idxVal models.Expression
		idxVal, err = ParseExpression(toks)
		if err != nil {
			return nil, err
		}
		if idxVal != nil {
			idx2 = &idxVal
		}

		return &ArraySliceExpression{
			Array: arr,
			Begin: idx,
			End:   idx2,
			loc:   arr.SourceLocation(),
		}, nil
	}

	if idx == nil {
		return nil, &models.InterpreterError{
			Message:        "expected array index expression",
			SourceLocation: beginLoc,
		}
	}

	return &ArrayAccessExpression{
		Array: arr,
		Index: *idx,
		loc:   arr.SourceLocation(),
	}, nil
}
