package codec

type Cloneable interface{ Clone() interface{} }

func Clone(v interface{}) interface{} {
	if v != nil {
		return v.(Cloneable).Clone()
	}
	return v
}
