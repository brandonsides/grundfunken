package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/brandonksides/grundfunken/models"
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
	"len": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)
			return len(list), nil
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
	"prepend": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[1].([]any)
			return append([]any{args[0]}, list...), nil
		},
	},
	"append": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list := args[0].([]any)

			newList := make([]any, len(list))
			copy(newList, list)

			return append(newList, args[1]), nil
		},
	},
	"concat": &BuiltinFunction{
		Argc: 2,
		Fn: func(args []any) (any, *models.InterpreterError) {
			list1 := args[0].([]any)
			newList := make([]any, len(list1))
			copy(newList, list1)

			list2 := args[1].([]any)
			return append(newList, list2...), nil
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
	"input": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			fmt.Print(args[0])
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("could not read input"),
					SourceLocation: models.SourceLocation{},
				}
			}
			return input, nil
		},
	},
	"print": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			fmt.Println(args[0])
			return nil, nil
		},
	},
	"sleep": &BuiltinFunction{
		Argc: 1,
		Fn: func(a []any) (any, *models.InterpreterError) {
			t := a[0].(int)
			time.Sleep(time.Duration(t) * time.Millisecond)
			return t, nil
		},
	},
	"parseInt": &BuiltinFunction{
		Argc: 1,
		Fn: func(args []any) (any, *models.InterpreterError) {
			str := args[0].(string)
			var num int
			_, err := fmt.Sscanf(str, "%d", &num)
			if err != nil {
				return nil, &models.InterpreterError{
					Err:            fmt.Errorf("could not parse int from string"),
					SourceLocation: models.SourceLocation{},
				}
			}
			return num, nil
		},
	},
}
