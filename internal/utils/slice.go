package utils

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// note: if subSlice is empty it will return true, this is mostly used for verifying if the scopes requested by user are already assigned to user
func ContainsSlice[T comparable](slice []T, subSlice []T) bool {
	for _, elem := range subSlice {
		if !Contains(slice, elem) {
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
