package utils

func Any(args ...bool) bool {
	for _, arg := range args {
		if !arg {
			return false
		}
	}
	return true
}