package stuff

// InArray returns true if the element is in the array
func InArray(str string, strings []string) bool {
	for _, s := range strings {
		if s == str {
			return true
		}
	}
	return false
}
