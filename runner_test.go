package gofuncount

import (
	"reflect"
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
		want    []CountInfo
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				root: "testdata/src",
				conf: &Config{
					IncludeTests: true,
				},
			},
			want: []CountInfo{
				{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.root, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
