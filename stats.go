package gofuncount

import (
	"math"
	"sort"
)

type Data struct {
	floats []float64
}

func NewData[T float64 | int | uint](values []T) Data {
	var floats []float64

	for _, v := range values {
		floats = append(floats, float64(v))
	}

	sort.Float64s(floats)

	return Data{
		floats: floats,
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
