package types

type ObjectType struct {
	Fields map[string]Type
}

func (ot ObjectType) String() string {
	str := "{"
	first := true
	for k, ty := range ot.Fields {
		if first {
			first = false
		} else {
			str += ", "
		}
		str += k + ": " + ty.String()
	}
	str += "}"
	return str
}

func Object(fields map[string]Type) ObjectType {
	return ObjectType{Fields: fields}
}
