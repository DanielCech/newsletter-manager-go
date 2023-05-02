package util

// Generic map
func MapFunc[T any, S any](data []T, f func(T) S) []S {
	mapped := make([]S, len(data))

	for i, e := range data {
		mapped[i] = f(e)
	}

	return mapped
}
