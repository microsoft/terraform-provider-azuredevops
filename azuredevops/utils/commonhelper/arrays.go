package commonhelper

// SelectArrayRange creates a new array from s, based on the provided start and end index
func SelectArrayRange(s []string, startIndex int, endIndex int) (ret []string) {
	for i, v := range s {
		if i >= startIndex && i <= endIndex {
			ret = append(ret, v)
		}
	}
	return
}
