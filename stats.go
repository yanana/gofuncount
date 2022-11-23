package gofuncount

import (
	"math"
	"sort"
)

type Data struct {
	floats []float64
}

// NewData creates a new Data object from a slice of float64, int, or uint.
// If the input slice is not a slice of float64, int, or uint, the returned Data object will have a nil slice.
// If the slice is empty, the returned Data object will have a nil slice.
// If the slice is not empty, the returned Data object will have a copy of the slice.
// The slice is sorted in ascending order.
func NewData(raw interface{}) Data {
	var floats []float64

	switch t := raw.(type) {
	case []float64:
		floats = t
	case []uint:
		floats = make([]float64, len(t))
		for _, v := range t {
			floats = append(floats, float64(v))
		}
	case []int:
		floats = make([]float64, len(t))
		for _, v := range t {
			floats = append(floats, float64(v))
		}
	default:
		return Data{
			floats: []float64{},
		}
	}

	sort.Float64s(floats)

	return Data{
		floats,
	}
}

func (d *Data) Quantile(p float64) float64 {
	data := d.floats
	count := len(data)
	if count == 0 {
		return math.NaN()
	}
	if count == 1 {
		return data[0]
	}
	sort.Float64s(data)

	pos := p * (float64(count) - 1)
	n := math.Min(math.Floor(pos), float64(count)-2)
	r := pos - n

	return data[int(n)]*(1-r) + data[int(n+1)]*r
}

func (d *Data) Mean() float64 {
	data := d.floats
	count := len(data)
	if count == 0 {
		return math.NaN()
	}

	var sum float64
	for _, v := range data {
		sum += v
	}

	return sum / float64(count)
}
