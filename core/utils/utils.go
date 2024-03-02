package utils

func Any(args ...bool) bool {
	for _, arg := range args {
		if !arg {
			return false
		}
	}
	return true
}

func MergeSortedArrays(a, b []string) []string {
	merged := make([]string, len(a)+len(b))
	var aIndex, bIndex int

	for i := 0; i < len(a)+len(b); i++ {
		// TODO: removed code duplicate, make it clearer
		if aIndex < len(a) && bIndex < len(b) {
			if len(a[aIndex]) >= len(b[bIndex]) {
				merged[i] = a[aIndex]
				aIndex++
			} else {
				merged[i] = b[bIndex]
				bIndex++
			}
			continue
		}
		if aIndex < len(a) {
			merged[i] = a[aIndex]
			aIndex++
		} else {
			merged[i] = b[bIndex]
			bIndex++
		}
	}
	return merged
}
