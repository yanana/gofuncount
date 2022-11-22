package gofuncount

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestData(t *testing.T) {
	t.Parallel()

	t.Run("Quantile", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			name   string
			data   Data
			p      float64
			want   float64
			assert assert.ComparisonAssertionFunc
		}{
			{
				name: "empty",
				data: NewData([]int{}),
				p:    0.5,
				want: math.NaN(),
				assert: func(t assert.TestingT, want, got interface{}, msgAndArgs ...interface{}) bool {
					return assert.True(t, math.IsNaN(got.(float64)), msgAndArgs...)
				},
			},
			{
				name:   "single",
				data:   NewData([]float64{99.9}),
				p:      0.5,
				want:   99.9,
				assert: assert.Equal,
			},
			{
				name:   "ten/0.5",
				data:   NewData([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				p:      0.5,
				want:   5.5,
				assert: assert.Equal,
			},
			{
				name:   "ten/0.9",
				data:   NewData([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				p:      0.9,
				want:   9.1,
				assert: assert.Equal,
			},
			{
				name: "ten/0.95",
				data: NewData([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				p:    0.95,
				want: 9.55,
				assert: func(t assert.TestingT, want, got interface{}, msgAndArgs ...interface{}) bool {
					return assert.InDelta(t, want, got, 0.1, msgAndArgs...)
				},
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				tc.assert(t, tc.want, tc.data.Quantile(tc.p))
			})
		}
	})

}
