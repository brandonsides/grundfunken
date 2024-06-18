package main

import (
	"fmt"

	"github.com/brandonksides/phonk/models"
)

type BuiltinFunction struct {
	Argc int
	Fn   func([]any) (any, *models.InterpreterError)
}

func (f BuiltinFunction) Call(args []any) (ret any, err *models.InterpreterError) {
	defer func() {
		if r := recover(); r != nil {
			err = &models.InterpreterError{
				Err: fmt.Errorf("panic: %v", r),
			}
		}
	}()

	if len(args) > f.Argc {
		return nil, &models.InterpreterError{
			Err: fmt.Errorf("expected %d arguments, got %d", f.Argc, len(args)),
		}
	}
	return f.Fn(args)
}

var builtins = map[string]any{
	"at": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)
			index := args[1].(int)
			if index < 0 || index >= len(list) {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("index out of bounds"),
					SourceLocation: models.SourceLocation{},
				}
			}
			return list[index], nil
		},
	},
	"len": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)
			return len(list), nil
		},
	},
	"slice": &BuiltinFunction{
		Argc: 3,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)
			start := args[1].(int)
			end := args[2].(int)
			if start < 0 {
				start = len(list) + start + 1
				if start < 0 {
					return nil, &models.InterpreterError{
						Err:            fmt.Errorf("start index out of bounds"),
						SourceLocation: models.SourceLocation{},
					}
				}
			}
			if end < 0 {
				end = len(list) + end + 1
				if end < 0 {
					return nil, &models.InterpreterError{
						Err:            fmt.Errorf("end index out of bounds"),
						SourceLocation: models.SourceLocation{},
					}
				}
			}

			return list[start:end], nil
		},
	},
	"range": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			start := args[0].(int)
			end := args[1].(int)

			ret := make([]any, 0, end-start)
			for start < end {
				ret = append(ret, start)
				start++
			}

			return ret, nil
		},
	},
	"lessThan": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v1, ok := args[0].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}
			v2, ok := args[1].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[1]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return v1 < v2, nil
		},
	},
	"greaterThan": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v1, ok := args[0].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}
			v2, ok := args[1].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[1]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return v1 > v2, nil
		},
	},
	"equals": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			return args[0] == args[1], nil
		},
	},
	"prepend": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[1].([]any)
			return append([]any{args[0]}, list...), nil
		},
	},
	"or": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v1, ok := args[0].(bool)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected bool, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}

			// short circuit
			if v1 {
				return true, nil
			}

			v2, ok := args[1].(bool)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected bool, got %T", args[1]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return v2, nil
		},
	},
	"and": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v1, ok := args[0].(bool)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected bool, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}

			// short circuit
			if !v1 {
				return false, nil
			}

			v2, ok := args[1].(bool)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected bool, got %T", args[1]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return v2, nil
		},
	},
	"not": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v, ok := args[0].(bool)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected bool, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return !v, nil
		},
	},
	"mod": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			v1, ok := args[0].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[0]),
					SourceLocation: models.SourceLocation{},
				}
			}
			v2, ok := args[1].(int)
			if !ok {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("expected int, got %T", args[1]),
					SourceLocation: models.SourceLocation{},
				}
			}
			return v1 % v2, nil
		},
	},
	"append": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)
			return append(list, args[1]), nil
		},
	},
	"concat": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list1 := args[0].([]any)
			list2 := args[1].([]any)
			return append(list1, list2...), nil
		},
	},
	"concatStr": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			str1 := args[0].(string)
			str2 := args[1].(string)
			return str1 + str2, nil
		},
	},
	"atStr": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			str := args[0].(string)
			index := args[1].(int)
			if index < 0 || index >= len(str) {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("index out of bounds"),
					SourceLocation: models.SourceLocation{},
				}
			}
			return string(str[index]), nil
		},
	},
	"lenStr": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			str := args[0].(string)
			return len(str), nil
		},
	},
	"sliceStr": &BuiltinFunction{
		Argc: 3,
		Fn: func(args []any) (any, *models.InterpreterError) {
			str := args[0].(string)
			start := args[1].(int)
			end := args[2].(int)
			if start < 0 {
				start = len(str) + start + 1
				if start < 0 {
					return nil, &models.InterpreterError{
						Err:            fmt.Errorf("start index out of bounds"),
						SourceLocation: models.SourceLocation{},
					}
				}
			}
			if end < 0 {
				end = len(str) + end + 1
				if end < 0 {
					return nil, &models.InterpreterError{
						Err:            fmt.Errorf("end index out of bounds"),
						SourceLocation: models.SourceLocation{},
					}
				}
			}

			return str[start:end], nil
		},
	},
}
