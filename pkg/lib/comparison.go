package lib

func CompareSliceStr(slice1, slice2 []string) (hasSameLength, containsSameValue bool) {
	hasSameLength = len(slice1) == len(slice2)

	slice1ContainsSlice2 := true
loopSlice1:
	for _, item1 := range slice1 {
		_, slice1ContainsSlice2 = FindSlice(slice2, item1)
		if !slice1ContainsSlice2 {
			break loopSlice1
		}
	}

	slice2ContainsSlice1 := true
loopSlice2:
	for _, item2 := range slice2 {
		_, slice2ContainsSlice1 = FindSlice(slice1, item2)
		if !slice2ContainsSlice1 {
			break loopSlice2
		}
	}

	containsSameValue = slice1ContainsSlice2 && slice2ContainsSlice1

	return
}
