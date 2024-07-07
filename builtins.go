package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/brandonksides/grundfunken/models/types"
)

type BuiltinFunction struct {
	args []types.Arg
	ret  types.Type
	Fn   func([]any) (any, error)
}

func Builtin(args []types.Arg, ret types.Type, fn func([]any) (any, error)) types.Function {
	return &BuiltinFunction{
		args: args,
		ret:  ret,
		Fn:   fn,
	}
}

var _ types.Function = &BuiltinFunction{}

func (f BuiltinFunction) Call(args []any) (ret any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if len(args) > len(f.args) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(f.args), len(args))
	}
	return f.Fn(args)
}

func (f BuiltinFunction) Args() []types.Arg {
	return f.args
}

func (f BuiltinFunction) Return() types.Type {
	return f.ret
}

var builtins = map[string]any{
	"len": &BuiltinFunction{
		args: []types.Arg{{
			Name: "list",
			Type: types.List(types.PrimitiveTypeAny),
		}},
		ret: types.PrimitiveTypeInt,
		Fn: func(args []any) (any, error) {
			list := args[0].([]any)
			return len(list), nil
		},
	},
	"range": &BuiltinFunction{
		args: []types.Arg{
			{
				Name: "start",
				Type: types.PrimitiveTypeInt,
			},
			{
				Name: "end",
				Type: types.PrimitiveTypeInt,
			},
		},
		ret: types.List(types.PrimitiveTypeInt),
		Fn: func(args []any) (any, error) {
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
	"toString": &BuiltinFunction{
		args: []types.Arg{{
			Name: "val",
			Type: types.PrimitiveTypeAny,
		}},
		ret: types.PrimitiveTypeString,
		Fn: func(args []any) (any, error) {
			return fmt.Sprint(args[0]), nil
		},
	},
	"prepend": &BuiltinFunction{
		args: []types.Arg{{
			Name: "item",
			Type: types.PrimitiveTypeAny,
		}, {
			Name: "list",
			Type: types.List(types.PrimitiveTypeAny),
		}},
		ret: types.List(types.PrimitiveTypeAny),
		Fn: func(args []any) (any, error) {
			list := args[1].([]any)
			return append([]any{args[0]}, list...), nil
		},
	},
	"append": &BuiltinFunction{
		args: []types.Arg{{
			Name: "list",
			Type: types.List(types.PrimitiveTypeAny),
		}, {
			Name: "item",
			Type: types.PrimitiveTypeAny,
		}},
		ret: types.List(types.PrimitiveTypeAny),
		Fn: func(args []any) (any, error) {
			list := args[0].([]any)

			newList := make([]any, len(list))
			copy(newList, list)

			return append(newList, args[1]), nil
		},
	},
	"concat": &BuiltinFunction{
		args: []types.Arg{{
			Name: "list1",
			Type: types.List(types.PrimitiveTypeAny),
		}, {
			Name: "list2",
			Type: types.List(types.PrimitiveTypeAny),
		}},
		ret: types.List(types.PrimitiveTypeAny),
		Fn: func(args []any) (any, error) {
			list1 := args[0].([]any)
			newList := make([]any, len(list1))
			copy(newList, list1)

			list2 := args[1].([]any)
			return append(newList, list2...), nil
		},
	},
	"concatStr": &BuiltinFunction{
		args: []types.Arg{{
			Name: "str1",
			Type: types.PrimitiveTypeString,
		}, {
			Name: "str2",
			Type: types.PrimitiveTypeString,
		}},
		ret: types.PrimitiveTypeString,
		Fn: func(args []any) (any, error) {
			str1 := args[0].(string)
			str2 := args[1].(string)
			return str1 + str2, nil
		},
	},
	"atStr": &BuiltinFunction{
		args: []types.Arg{{
			Name: "str",
			Type: types.PrimitiveTypeString,
		}, {
			Name: "index",
			Type: types.PrimitiveTypeInt,
		}},
		ret: types.PrimitiveTypeString,
		Fn: func(args []any) (any, error) {
			str := args[0].(string)
			index := args[1].(int)
			if index < 0 || index >= len(str) {
				return nil, fmt.Errorf("index out of bounds (%d); len is %d", index, len(str))
			}
			return string(str[index]), nil
		},
	},
	"lenStr": &BuiltinFunction{
		args: []types.Arg{{
			Name: "str",
			Type: types.PrimitiveTypeString,
		}},
		ret: types.PrimitiveTypeInt,
		Fn: func(args []any) (any, error) {
			str := args[0].(string)
			return len(str), nil
		},
	},
	"sliceStr": &BuiltinFunction{
		args: []types.Arg{{
			Name: "str",
			Type: types.PrimitiveTypeString,
		}, {
			Name: "start",
			Type: types.PrimitiveTypeInt,
		}, {
			Name: "end",
			Type: types.PrimitiveTypeInt,
		}},
		ret: types.PrimitiveTypeString,
		Fn: func(args []any) (any, error) {
			str := args[0].(string)
			start := args[1].(int)
			end := args[2].(int)
			if start < 0 {
				start = len(str) + start + 1
				if start < 0 || start > len(str) {
					return nil, fmt.Errorf("start index out of bounds (%d); len is %d", start, len(str))
				}
			}
			if end < 0 {
				end = len(str) + end + 1
				if end < 0 || end > len(str) {
					return nil, fmt.Errorf("end index out of bounds (%d); len is %d", end, len(str))
				}
			}

			return str[start:end], nil
		},
	},
	"input": &BuiltinFunction{
		args: []types.Arg{{
			Name: "prompt",
			Type: types.PrimitiveTypeString,
		}},
		ret: types.PrimitiveTypeString,
		Fn: func(args []any) (any, error) {
			fmt.Print(args[0])
			reader := bufio.NewReader(os.Stdin)
			return reader.ReadString('\n')
		},
	},
	"print": &BuiltinFunction{
		args: []types.Arg{{
			Name: "val",
			Type: types.PrimitiveTypeAny,
		}},
		ret: types.PrimitiveTypeUnit,
		Fn: func(args []any) (any, error) {
			fmt.Println(args[0])
			return nil, nil
		},
	},
	"sleep": &BuiltinFunction{
		args: []types.Arg{{
			Name: "time",
			Type: types.PrimitiveTypeInt,
		}},
		ret: types.PrimitiveTypeUnit,
		Fn: func(a []any) (any, error) {
			t := a[0].(int)
			time.Sleep(time.Duration(t) * time.Millisecond)
			return t, nil
		},
	},
	"parseInt": &BuiltinFunction{
		args: []types.Arg{{
			Name: "str",
			Type: types.PrimitiveTypeString,
		}},
		ret: types.PrimitiveTypeInt,
		Fn: func(args []any) (any, error) {
			str := args[0].(string)
			var num int
			_, err := fmt.Sscanf(str, "%d", &num)
			if err != nil {
				return nil, fmt.Errorf("could not parse int from string \"%s\"", str)
			}
			return num, nil
		},
	},
}
