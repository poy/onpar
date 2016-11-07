package fibonacci

func Fibonacci(n int) int {
	if n <= 1 {
		return 1
	}

	return Fibonacci(n-2) + Fibonacci(n-1)
}
