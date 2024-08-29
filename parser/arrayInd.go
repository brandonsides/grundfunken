package parser

import (
	"fmt"

	"github.com/brandonksides/grundfunken/models"
	"github.com/brandonksides/grundfunken/models/expressions"
	"github.com/brandonksides/grundfunken/models/types"
	"github.com/brandonksides/grundfunken/tokens"
)

type ArrayAccessExpression struct {
	Array expressions.Expression
	Index expressions.Expression
	loc   *models.SourceLocation
}

func (aae *ArrayAccessExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	t, err := aae.Array.Type(tb)
	if err != nil {
		return nil, err
	}

	tList, ok := t.(types.ListType)
	if !ok {
		return nil, &models.InterpreterError{
			Message:        fmt.Sprintf("expected list; got %s", t),
			SourceLocation: aae.Array.SourceLocation(),
		}
	}

	return tList.ElementType, nil
}

func (aae *ArrayAccessExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
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

func (aae *ArrayAccessExpression) SourceLocation() *models.SourceLocation {
	ret := *aae.loc
	return &ret
}

type ArraySliceExpression struct {
	Array expressions.Expression
	Begin *expressions.Expression
	End   *expressions.Expression
	loc   *models.SourceLocation
}

func (ase *ArraySliceExpression) Type(tb types.TypeBindings) (types.Type, *models.InterpreterError) {
	t, err := ase.Array.Type(tb)
	if err != nil {
		return nil, err
	}

	if ase.Begin != nil {
		t, err := (*ase.Begin).Type(tb)
		if err != nil {
			return nil, err
		}
		if t != types.PrimitiveTypeInt {
			return nil, &models.InterpreterError{
				Message:        fmt.Sprintf("expected int; got %s", t),
				SourceLocation: (*ase.Begin).SourceLocation(),
			}
		}
	}

	if ase.End != nil {
		t, err := (*ase.End).Type(tb)
		if err != nil {
			return nil, err
		}
		if t != types.PrimitiveTypeInt {
			return nil, &models.InterpreterError{
				Message:        fmt.Sprintf("expected int; got %s", t),
				SourceLocation: (*ase.End).SourceLocation(),
			}
		}
	}

	return t, nil
}

func (ase *ArraySliceExpression) Evaluate(bindings expressions.Bindings) (any, *models.InterpreterError) {
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

func (ase *ArraySliceExpression) SourceLocation() *models.SourceLocation {
	ret := *ase.loc
	return &ret
}

func parseArrayIndex(arr expressions.Expression, toks *tokens.TokenStack) (exp expressions.Expression, err *models.InterpreterError) {
	beginLoc := toks.CurrentSourceLocation()

	tok, ok := toks.Peek()
	if !ok {
		return nil, &models.InterpreterError{
			Message:        "expected array index expression",
			SourceLocation: arr.SourceLocation(),
		}
	}

	var idx *expressions.Expression
	if tok.Type != tokens.COLON {
		var idxVal expressions.Expression
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

		var idx2 *expressions.Expression
		var idxVal expressions.Expression
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
