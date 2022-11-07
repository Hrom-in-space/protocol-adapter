package utils

import "strings"

// RemoveBeforeFirst - remove prefix before first separator.
func RemoveBeforeFirst(str, sep string) string {
	_, result, ok := strings.Cut(str, sep)
	if !ok {
		result = str
	}

	return result
}

// PackagePath - remove path before "pkgDir/pkgFile.go".
func PackagePath(path string) string {
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	return path[idx+1:]
}
