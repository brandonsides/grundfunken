package types

type ListType struct {
	ElementType Type
}

func (lt ListType) String() string {
	return "[" + lt.ElementType.String() + "]"
}

func List(elementType Type) ListType {
	return ListType{ElementType: elementType}
}
