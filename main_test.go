package main

import (
	"reflect"
	"testing"
)

func Test_calculateTotalSize(t *testing.T) {
	type args struct {
		folderSizes map[string]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Calculate total size",
			args: args{
				folderSizes: map[string]int{
					"file1": 1024,
					"file2": 2048,
				},
			},
			want: 3072,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTotalSize(tt.args.folderSizes); got != tt.want {
				t.Errorf("calculateTotalSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatSize(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Format size bytes",
			args: args{
				size: 512,
			},
			want: "512 bytes",
		},
		{
			name: "Format size KB",
			args: args{
				size: 1024,
			},
			want: "1 KB",
		},
		{
			name: "Format size MB",
			args: args{
				size: 1024 * 1024,
			},
			want: "1.00 MB",
		},
		{
			name: "Format size GB",
			args: args{
				size: 1024 * 1024 * 1024,
			},
			want: "1.00 GB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatSize(tt.args.size); got != tt.want {
				t.Errorf("formatSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectBestFolders(t *testing.T) {
	type args struct {
		folderSizes map[string]int
		totalSize   int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "Select best folders",
			args: args{
				folderSizes: map[string]int{
					"file1": 1024 * 1024,
					"file2": 2048 * 1024,
				},
				totalSize: 3072 * 1024,
			},
			want: map[string]int{
				"file1": 1024 * 1024,
				"file2": 2048 * 1024,
			},
		},
		{
			name: "Select best folders",
			args: args{
				folderSizes: map[string]int{
					"file1": 1024,
					"file2": 4048,
					"file3": 2048,
				},
				totalSize: 5072,
			},
			want: map[string]int{
				"file1": 1024,
				"file2": 4048,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectBestFolders(tt.args.folderSizes, tt.args.totalSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectBestFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateSize(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Calculate file size",
			args: args{
				source: "testing/file.txt",
			},
			want: 16,
		},
		{
			name: "Calculate directory size",
			args: args{
				source: "testing",
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateSize(tt.args.source); got != tt.want {
				t.Errorf("calculateSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSizeInBytes(t *testing.T) {
	type args struct {
		sizeStr string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Get size in bytes",
			args: args{
				sizeStr: "1024",
			},
			want: 1024,
		},
		{
			name: "Get size in bytes",
			args: args{
				sizeStr: "1024KB",
			},
			want: 1024 * 1024,
		},
		{
			name: "Get size in bytes",
			args: args{
				sizeStr: "1024MB",
			},
			want: 1024 * 1024 * 1024,
		},
		{
			name: "Get size in bytes",
			args: args{
				sizeStr: "1.5GB",
			},
			want: 1024 * 1024 * 1024 * 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSizeInBytes(tt.args.sizeStr); got != tt.want {
				t.Errorf("getSizeInBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
