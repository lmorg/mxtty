package elementTable

import (
	"strings"
)

func stringToInterfaceTrim(s []string, max int) []interface{} {
	slice := make([]interface{}, max)

	if max <= len(s) {
		var i int
		for ; i < max; i++ {
			slice[i] = s[i]
		}

		return slice
	}

	var i int
	for ; i < len(s); i++ {
		slice[i] = s[i]
	}

	for ; i < max; i++ {
		slice[i] = ""
	}

	return slice
}

func stringToInterfaceMerge(s []string, max int) []interface{} {
	slice := make([]interface{}, max)

	switch {
	case max == 0:
		// return empty slice

	case max < len(s):
		var i int
		for ; i < max-1; i++ {
			slice[i] = s[i]
		}
		slice[i] = strings.Join(s[i:], " ")

	case max == len(s):
		var i int
		for ; i < max; i++ {
			slice[i] = s[i]
		}

	case max > len(s):
		var i int
		for ; i < len(s); i++ {
			slice[i] = s[i]
		}
		for ; i < max; i++ {
			slice[i] = ""
		}
	}

	return slice
}
