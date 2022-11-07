package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"protocol-adapter/internal/utils"
)

func TestRemoveBeforeFirst(t *testing.T) {
	type args struct {
		str string
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty string", args{str: "", sep: "."}, ""},
		{"No separator", args{str: "abc", sep: "."}, "abc"},
		{"Split", args{str: "abc.def", sep: "."}, "def"},
		{"Double separator", args{str: "ab.cd.ed", sep: "."}, "cd.ed"},
		{"String of separators", args{str: "...", sep: "."}, ".."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, utils.RemoveBeforeFirst(tt.args.str, tt.args.sep), tt.want)
		})
	}
}

func TestPackagePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty", args{""}, ""},
		{"no path", args{"abc"}, "abc"},
		{"Short path", args{"/file"}, "/file"},
		{"Long path", args{"project/pkg/file"}, "pkg/file"},
		{"Normal absolute path", args{"/pkg/file"}, "pkg/file"},
		{"Normal relative path", args{"pkg/file"}, "pkg/file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, utils.PackagePath(tt.args.path), "PackagePath(%v)", tt.args.path)
		})
	}
}
