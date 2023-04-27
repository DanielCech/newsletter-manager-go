package util

import "fmt"

func Ptr[T any](v T) *T {
	return &v
}

func EqualPtrValue[T comparable](expected, actual *T) bool {
	return expected == actual || (expected != nil && actual != nil && *expected == *actual)
}

func PtrWithoutError[T any](v T, err error) *T {
	if err != nil {
		panic(fmt.Errorf("getting value: %w", err))
	}
	return &v
}

func ValueWithoutError[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Errorf("getting value: %w", err))
	}
	return v
}
