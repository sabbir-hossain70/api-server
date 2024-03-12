package math

// this is for average
func Average(xs []float64) float64 {
	total := float64(0)
	for _, x := range xs {
		total += x
	}
	return total / float64(len(xs))
}
func Min(xs []float64) float64 {
	min := xs[0]
	for _, z := range xs {
		if z < min {
			min = z
		}
	}
	return min
}
