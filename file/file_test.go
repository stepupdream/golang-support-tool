package file

import (
	"reflect"
	"testing"
)

func TestBaseNamesByArray(t *testing.T) {
	type args struct {
		paths         []string
		withExtension bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "BaseNamesByArray1",
			args: args{
				paths:         []string{"C:/develop/aaa.csv", "C:/develop/bbb.csv"},
				withExtension: false,
			},
			want: []string{"aaa", "bbb"},
		},
		{
			name: "BaseNamesByArray2",
			args: args{
				paths:         []string{"C:/develop/aaa.csv", "C:/develop/bbb.csv"},
				withExtension: true,
			},
			want: []string{"aaa.csv", "bbb.csv"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseNamesByArray(tt.args.paths, tt.args.withExtension); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseNamesByArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
