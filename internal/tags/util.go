package tags

import (
	"strings"
)

// Splits val and trims each item
func SplitTrim(val string, split string) []string {
	if val == "" {
		return nil
	}

	out := strings.Split(val, split)
	for index, item := range out {
		out[index] = strings.TrimSpace(item)
	}

	return out
}

// Checks the equality of 2 slices
func CompareArr[T comparable](arr1 []T, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
