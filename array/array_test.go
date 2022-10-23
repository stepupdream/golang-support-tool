package array

import (
	"reflect"
	"testing"
)

func TestSliceString(t *testing.T) {
	type args struct {
		all   []string
		start string
		end   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "max",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "max",
			},
			want: []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
		},
		{
			name: "max2",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_1_0",
				end:   "max",
			},
			want: []string{"1_0_1_0", "1_0_2_0"},
		},
		{
			name: "target",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "1_0_1_0",
			},
			want: []string{"1_0_0_0", "1_0_1_0"},
		},
		{
			name: "target2",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "1_0_2_0",
			},
			want: []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
		},
		{
			name: "startEmpty",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "",
				end:   "1_0_0_0",
			},
			want: []string{"1_0_0_0"},
		},
		{
			name: "next",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_1_0",
				end:   "next",
			},
			want: []string{"1_0_1_0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceString(tt.args.all, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextArrayValue(t *testing.T) {
	type args struct {
		allValues []string
		nowValue  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				allValues: []string{"a", "b", "c", "d"},
				nowValue:  "b",
			},
			want: "c",
		},
		{
			name: "test2",
			args: args{
				allValues: []string{"a", "b", "c", "d"},
				nowValue:  "d",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextArrayValue(tt.args.allValues, tt.args.nowValue); got != tt.want {
				t.Errorf("NextArrayValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringUnique(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "StringUnique",
			args: args{[]string{"a", "e", "b", "e", "d"}},
			want: []string{"a", "e", "b", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringUnique(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntUnique(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "IntUnique",
			args: args{[]int{1, 1, 2, 2, 3, 5}},
			want: []int{1, 2, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntUnique(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckStringByIndex(t *testing.T) {
	type args struct {
		rows  [][]string
		index int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "PluckStringByIndex",
			args: args{
				rows:  [][]string{{"a", "aaa"}, {"b", "bbb"}, {"c", "ccc"}},
				index: 0,
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluckStringByIndex(tt.args.rows, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckStringByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
