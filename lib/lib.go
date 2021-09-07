package lib

func IntArrayEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func EraseFromSlice(sl []int, v int) []int {
	for i, a := range sl {
		if a == v {
			return Remove(sl, i)
		}
	}
	return sl
}

func EraseMultiFromSlice(sl []int, v []int) []int {
	for i, a := range sl {
		for e := range v {
			if a == e {
				return Remove(sl, i)
			}
		}

	}
	return sl
}

func Remove(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
