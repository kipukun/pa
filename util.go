package pa

func sum[T Number](a []T) T {
	var sum T
	for _, val := range a {
		sum += val
	}
	return sum
}

func mean[T Number](a []T) T {
	var sum, count T
	for _, val := range a {
		sum += val
		count++
	}
	return sum / count
}
