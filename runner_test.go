package gofuncount

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		root string
		conf *Config
	}
	tests := []struct {
		name    string
		args    args
		want    Counts
		wantErr bool
	}{
		{
			name: "try to parse a non-existent file",
			args: args{
				root: "testdata/src/x/y.go",
				conf: &Config{
					IncludeTests: true,
				},
			},
			wantErr: true,
		},
		{
			name: "parse a file",
			args: args{
				root: "testdata/src/x/x.go",
				conf: &Config{
					IncludeTests: true,
				},
			},
			want: map[string][]*CountInfo{
				"testdata/src/x": {{
					Package:  "main",
					Name:     "init",
					FileName: "testdata/src/x/x.go",
					StartsAt: 5,
					EndsAt:   7,
				},
					{
						Package:  "main",
						Name:     "main",
						FileName: "testdata/src/x/x.go",
						StartsAt: 9,
						EndsAt:   12,
					},
				},
			},
		},
		{
			name: "parse files in a directory",
			args: args{
				root: "testdata/src/x",
				conf: &Config{
					IncludeTests: true,
				},
			},
			want: map[string][]*CountInfo{
				"testdata/src/x": {{
					Package:  "main",
					Name:     "init",
					FileName: "testdata/src/x/x.go",
					StartsAt: 5,
					EndsAt:   7,
				},
					{
						Package:  "main",
						Name:     "main",
						FileName: "testdata/src/x/x.go",
						StartsAt: 9,
						EndsAt:   12,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &Runner{
				Conf: tt.args.conf,
			}
			got, err := runner.Run(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, got, tt.want)
		})
	}
}

var epsilon = math.Nextafter(1, 2) - 1

func StatsEquals(t assert.TestingT, expected, actual interface{}, msgs ...interface{}) bool {
	if t, ok := t.(*testing.T); ok {
		t.Helper()
	}

	ex, ok := expected.(map[string]*Stats)
	if !ok {
		return assert.Fail(t, "expected is not a map[string]*Stats", msgs...)
	}

	ac, ok := actual.(map[string]*Stats)
	if !ok {
		return assert.Fail(t, "actual is not a map[string]*Stats", msgs...)
	}

	if len(ex) != len(ac) {
		t.Errorf("expected %d packages, got %d", len(ex), len(ac))
		return false
	}

	for k, v := range ex {
		actualStats, ok := ac[k]
		if !ok {
			t.Errorf("expected package %s, got none", k)
			return false
		}

		if statsEquals(t, v, actualStats) {
			return false
		}
	}

	return true
}

func statsEquals(t assert.TestingT, expected, actual *Stats) bool {
	if t, ok := t.(*testing.T); ok {
		t.Helper()
	}

	if !assert.InEpsilon(t, expected.MeanLines, actual.MeanLines, epsilon) {
		return false
	}
	if !assert.InEpsilon(t, expected.MedianLines, actual.MedianLines, epsilon) {
		return false
	}
	if !assert.InEpsilon(t, expected.NinetyFivePercentileLines, actual.NinetyFivePercentileLines, epsilon) {
		return false
	}
	if !assert.InEpsilon(t, expected.NinetyNinePercentileLines, actual.NinetyNinePercentileLines, epsilon) {
		return false
	}

	return true
}

func TestCounts(t *testing.T) {
	t.Parallel()

	t.Run("Stats()", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			cs     Counts
			want   map[string]*Stats
			assert assert.ComparisonAssertionFunc
		}{
			{
				name: "success",
				cs: Counts{"main": []*CountInfo{
					{
						Package:  "main",
						Name:     "init",
						FileName: "testdata/src/x/x.go",
						StartsAt: 5,
						EndsAt:   7,
					},
					{
						Package:  "main",
						Name:     "foo",
						FileName: "testdata/src/x/x.go",
						StartsAt: 8,
						EndsAt:   9,
					},
					{
						Package:  "main",
						Name:     "bar",
						FileName: "testdata/src/x/x.go",
						StartsAt: 10,
						EndsAt:   13,
					},
					{
						Package:  "main",
						Name:     "bar",
						FileName: "testdata/src/x/x.go",
						StartsAt: 10,
						EndsAt:   14,
					},
				},
				},
				want: map[string]*Stats{
					"main": {
						MeanLines:                 2.5,
						MedianLines:               2.5,
						NinetyFivePercentileLines: 3.85,
						NinetyNinePercentileLines: 3.97,
					},
				},
				assert: StatsEquals,
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				tc.assert(t, tc.want, tc.cs.Stats())
			})
		}
	})
}
