package codec

func Clone(v interface{}) interface{} {
	if v != nil {
		// type Cloneable interface{ Clone() interface{} }
		return v.(interface{ Clone() interface{} }).Clone()
	}
	return v
}
