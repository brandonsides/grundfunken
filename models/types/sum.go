package types

type sumType struct {
	Types []Type
}

func (t sumType) String() string {
	str := ""
	for i, ty := range t.Types {
		if i > 0 {
			str += " | "
		}
		str += ty.String()
	}
	return str
}

func Sum(types ...Type) Type {
	ret := sumType{Types: make([]Type, 0, len(types))}
	for _, t := range types {
		retSuper, err := IsSuperTo(ret, t)
		if err != nil {
			return nil
		}
		if tAsSum, ok := t.(sumType); ok {
			flattenedT := Sum(tAsSum.Types...)

			if flattenedTSum, ok := flattenedT.(sumType); ok {
				ret.Types = append(ret.Types, flattenedTSum.Types...)
			} else {
				ret.Types = append(ret.Types, flattenedT)
			}
		} else if retSuper {
			continue
		} else {
			for i := 0; i < len(ret.Types); i++ {
				tSuper, err := IsSuperTo(t, ret.Types[i])
				if err != nil {
					return nil
				}
				if tSuper {
					ret = sumType{Types: append(append([]Type{}, ret.Types[:i]...), ret.Types[i+1:]...)}
				} else {
					i++
				}
			}
			ret.Types = append(ret.Types, t)
		}
	}

	if len(ret.Types) == 1 {
		return ret.Types[0]
	}

	return ret
}
