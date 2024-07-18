package types

type PrimitiveType uint8

const (
	PrimitiveTypeInt PrimitiveType = iota
	PrimitiveTypeString
	PrimitiveTypeBool
	PrimitiveTypeUnit
	PrimitiveTypeAny
)

func (t PrimitiveType) String() string {
	switch t {
	case PrimitiveTypeInt:
		return "int"
	case PrimitiveTypeString:
		return "string"
	case PrimitiveTypeBool:
		return "bool"
	case PrimitiveTypeAny:
		return "any"
	case PrimitiveTypeUnit:
		return "unit"
	default:
		return "unknown"
	}
}

func ParsePrimitive(s string) PrimitiveType {
	switch s {
	case "int":
		return PrimitiveTypeInt
	case "string":
		return PrimitiveTypeString
	case "bool":
		return PrimitiveTypeBool
	case "unit":
		return PrimitiveTypeUnit
	case "any":
		return PrimitiveTypeAny
	default:
		return PrimitiveTypeAny
	}
}
