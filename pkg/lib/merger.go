package lib

// Merge a struct to another struct
func Merge(from interface{}, to interface{}) error {
	j, e := JSONMarshal(from)
	if nil == e {
		e = JSONUnmarshal(j, to)
	}

	return e
}
