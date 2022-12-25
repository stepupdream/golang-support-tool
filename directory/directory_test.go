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

func TestMaxFileName(t *testing.T) {
	type args struct {
		directoryPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "MaxFileName",
			args: args{
				directoryPath: "../directory/test",
			},
			want: "1_0_1_0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFileName(tt.args.directoryPath); got != tt.want {
				t.Errorf("MaxFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
