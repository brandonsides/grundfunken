package types

import (
	"fmt"
)

type TypeBindings map[string]Type

type Type interface {
	fmt.Stringer
}

func TypeOf(v interface{}) (Type, error) {
	switch v := v.(type) {
	case nil:
		return PrimitiveTypeUnit, nil
	case int:
		return PrimitiveTypeInt, nil
	case string:
		return PrimitiveTypeString, nil
	case bool:
		return PrimitiveTypeBool, nil
	case []interface{}:
		elemTypes := make([]Type, 0)
		for _, elem := range v {
			typ, err := TypeOf(elem)
			if err != nil {
				return nil, err
			}
			elemTypes = append(elemTypes, typ)
		}
		return List(Sum(elemTypes...)), nil
	case map[string]interface{}:
		fieldTypes := make(map[string]Type)
		for k, v := range v {
			var err error
			fieldTypes[k], err = TypeOf(v)
			if err != nil {
				return nil, err
			}
		}
		return Object(fieldTypes), nil
	case Function:
		typs := make([]Type, 0)
		for _, arg := range v.Args() {
			typs = append(typs, arg.Type)
		}
		return Func(typs, v.Return()), nil
	default:
		return nil, fmt.Errorf("unknown type %T", v)
	}
}

type Function interface {
	Call([]any) (any, error)
	Args() []Arg
	Return() Type
}

type Arg struct {
	Name string
	Type Type
}

func IsSuperTo(t1, t2 Type) (bool, error) {
	if t2Sum, ok := t2.(sumType); ok {
		for _, t2Addend := range t2Sum.Types {
			superToAddend, err := IsSuperTo(t1, t2Addend)
			if err != nil {
				return false, err
			}
			if !superToAddend {
				return false, nil
			}
		}
	} else if t1Sum, ok := t1.(sumType); ok {
		for _, t1Addend := range t1Sum.Types {
			superToAddend, err := IsSuperTo(t1Addend, t2)
			if err != nil {
				return false, err
			}
			if superToAddend {
				return true, nil
			}
		}
		return false, nil
	}

	switch t1 := t1.(type) {
	case PrimitiveType:
		if t1 == PrimitiveTypeAny {
			return true, nil
		}
		if t2Prim, ok := t2.(PrimitiveType); ok {
			return t1 == t2Prim, nil
		}
		return false, nil
	case ListType:
		if t2List, ok := t2.(ListType); ok {
			return IsSuperTo(t1.ElementType, t2List.ElementType)
		}
		return false, nil
	case ObjectType:
		if t2Obj, ok := t2.(ObjectType); ok {
			for k, v1 := range t1.Fields {
				super, err := IsSuperTo(v1, t2Obj.Fields[k])
				if err != nil {
					return false, err
				}
				if _, ok := t2Obj.Fields[k]; !ok || !super {
					return false, nil
				}
			}
			return true, nil
		}
		return false, nil
	case FuncType:
		if t2Func, ok := t2.(FuncType); ok {
			if len(t1.ArgTypes) != len(t2Func.ArgTypes) {
				return false, nil
			}
			for i, arg := range t1.ArgTypes {
				super, err := IsSuperTo(t2Func.ArgTypes[i], arg)
				if err != nil {
					return false, err
				}
				if !super {
					return false, nil
				}
			}
			return IsSuperTo(t1.ReturnType, t2Func.ReturnType)
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown type %T", t1)
	}
}
