package utils

func TernaryOperator[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
