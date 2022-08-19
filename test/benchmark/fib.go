package benchmark

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func fib2(n int) int {
	if n < 2 {
		return n
	}
	var f1, f2, f3 = 0, 0, 1
	for i := 2; i <= n; i++ {
		f1 = f2
		f2 = f3
		f3 = f1 + f2
	}
	return f3
}

func fib3(n int) int {
	dp := make([]int, n+1)
	dp[1] = 1
	dp[2] = 1
	for i := 2; i < n; i++ {
		dp[i] = dp[i-2] + dp[i-1]
	}
	return dp[n]
}
