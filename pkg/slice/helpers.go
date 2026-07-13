package slice

import (
	"slices"
	"strconv"
)

// Contains Check item in slice string type
func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// ContainsRune Check item in map rune type
func ContainsRune(set map[rune]struct{}, item rune) bool {
	_, ok := set[item]
	return ok
}

// ContainsValue check if value exists in slice, no matter its type
func ContainsValue(slice []any, value any) bool {
	return slices.Contains(slice, value)
}

// IntToString Convert slice int to slice string
func IntToString(intSl []int) (str []string) {
	for i := range intSl {
		number := intSl[i]
		text := strconv.Itoa(number)
		str = append(str, text)
	}
	return str
}
