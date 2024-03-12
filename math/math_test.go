package math

import "testing"

type testpair struct {
	values []float64
	avg    float64
}

var tests = []testpair{
	{[]float64{1, 2, 3}, 2},
	{[]float64{33, 33, 33, 33}, 33},
	{[]float64{1, 2}, 1.5},
	{[]float64{1}, 1},
}

func TestAverage(t *testing.T) {
	for _, pair := range tests {
		v := pair.values
		a := pair.avg
		if Average(v) != a {
			t.Error(
				"For ", v, " Expected ", a, " found ", Average(v),
			)
		}
	}
}
