package shared

func Contains[T comparable](s []T, val T) bool {
	for _, el := range s {
		if el == val {
			return true
		}
	}
	return false
}
