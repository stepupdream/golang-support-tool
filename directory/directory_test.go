package directory

import "testing"

func TestExistMulti(t *testing.T) {
	type args struct {
		parentPaths []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ExistMulti",
			args: args{
				parentPaths: []string{"../directory", "../excel"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExistMulti(tt.args.parentPaths); got != tt.want {
				t.Errorf("ExistMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}
