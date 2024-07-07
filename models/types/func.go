package types

type FuncType struct {
	ArgTypes   []Type
	ReturnType Type
}

func (ft FuncType) String() string {
	return "func"
}

func Func(argTypes []Type, returnType Type) FuncType {
	return FuncType{ArgTypes: argTypes, ReturnType: returnType}
}
