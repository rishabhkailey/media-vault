package utils

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func ContainsSlice[T comparable](slice []T, subSlice []T) bool {
	var elementsMap map[T]bool = make(map[T]bool)
	for _, elem := range slice {
		elementsMap[elem] = true
	}
	for _, elem := range subSlice {
		if contains, ok := elementsMap[elem]; !contains || !ok {
			return false
		}
	}
	return true
}

func SliceMap[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
