package fileutils

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "can replace string invalid for filename", args: args{fileName: "a<b>c:d\\e\"f/g\\h|i?j*k"}, want: "abcdefghijk"},
		{name: "casing and whitespace preserved", args: args{fileName: "aB Cd"}, want: "aB Cd"},
		{name: "harmless symbol remain valid", args: args{fileName: "~!@#$%^&()[].,"}, want: "~!@#$%^&()[].,"},
		{name: "can normalize multiple whitespace", args: args{fileName: "this  is multiple    whitespace"}, want: "this is multiple whitespace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeFilename(tt.args.fileName); got != tt.want {
				t.Errorf("SanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasePath(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	tests := []struct {
		name string
		want string
	}{
		{name: "get path", want: filepath.Join(basepath, "../../")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BasePath(); got != tt.want {
				t.Errorf("BasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourcePath(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	tests := []struct {
		name string
		want string
	}{
		{name: "get path", want: filepath.Join(basepath, "../../", "resources")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResourcePath(); got != tt.want {
				t.Errorf("ResourcePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDataJson(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "no dir", args: args{path: ""}, want: []string{}, wantErr: true},
		{name: "no data json available", args: args{path: basepath}, want: []string{}, wantErr: false},
		{name: "scan resources for data json available", args: args{path: filepath.Join(basepath, "../../", "resources")}, want: []string{"1", "2"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDataJson(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !(len(got) == len(tt.want)) {
				t.Errorf("GetDataJson() got len = %v, want len %v", len(got), len(tt.want))
			}
		})
	}
}
