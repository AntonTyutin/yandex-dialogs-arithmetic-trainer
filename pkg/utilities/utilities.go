package utilities

import "math/rand"

func Coalesce[T comparable](values ...T) T {
	var empty T
	for _, value := range values {
		if value != empty {
			return value
		}
	}

	return empty
}

func Random[T any](values []T) T {
	return values[rand.Intn(len(values))]
}
