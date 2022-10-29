package csv

import (
	"reflect"
	"testing"
)

func TestDeleteCSV(t *testing.T) {
	type args struct {
		baseCSV     map[Key]string
		editCSV     map[Key]string
		filterNames []string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "DeleteBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
				editCSV: map[Key]string{
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
				filterNames: []string{"id", "name"},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteCSV(tt.args.baseCSV, tt.args.editCSV); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertCSV(t *testing.T) {
	type args struct {
		baseCSV map[Key]string
		editCSV map[Key]string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "InsertBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
				editCSV: map[Key]string{
					{Id: 4, Key: "id"}: "3", {Id: 3, Key: "name"}: "dddd",
				},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
				{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				{Id: 4, Key: "id"}: "3", {Id: 3, Key: "name"}: "dddd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertCSV(tt.args.baseCSV, tt.args.editCSV); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insertCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCSV(t *testing.T) {
	type args struct {
		baseCSV map[Key]string
		editCSV map[Key]string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "InsertBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
				},
				editCSV: map[Key]string{
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "eeee",
				},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "eeee",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateCSV(tt.args.baseCSV, tt.args.editCSV); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFileFirstContent(t *testing.T) {
	type args struct {
		directoryPath string
		fileName      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "LoadFileFirstContent",
			args: args{
				directoryPath: "test",
				fileName:      "sample.csv",
			},
			want: "sample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadFileFirstContent(tt.args.directoryPath, tt.args.fileName); got != tt.want {
				t.Errorf("LoadFileFirstContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadCsv(t *testing.T) {
	type args struct {
		filepath string
		isFilter bool
	}
	tests := []struct {
		name  string
		args  args
		want  [][]string
		want1 []string
	}{
		{
			name: "LoadCsv",
			args: args{
				filepath: "test/sample2.csv",
				isFilter: false,
			},
			want: [][]string{
				{"id", "sample", "#", "level"},
				{"1", "aaa", "2", "13"},
				{"2", "bbb", "3", "43"},
			},
			want1: []string{"id", "sample", "#", "level"},
		},
		{
			name: "LoadCsv",
			args: args{
				filepath: "test/sample2.csv",
				isFilter: true,
			},
			want: [][]string{
				{"id", "sample", "level"},
				{"1", "aaa", "13"},
				{"2", "bbb", "43"},
			},
			want1: []string{"id", "sample", "level"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadCsv(tt.args.filepath, tt.args.isFilter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadCsv() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoadCsv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
