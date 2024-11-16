package octal

import (
	"fmt"
	"strconv"
)

func Escape(b []byte) []byte {
	var escaped []byte

	for _, c := range b {
		//if c <= ' ' || c {
		escaped = append(escaped, []byte(fmt.Sprintf(`\%03o`, c))...)
		//	continue
		//}
		//escaped = append(escaped, b...)
	}

	return escaped
}

func Unescape(b []byte) []byte {
	var (
		c = make([]byte, len(b))
		j int
	)

	for i := 0; i < len(b); j++ {
		if b[i] != '\\' {
			c[j] = b[i]
			i++
			continue
		}

		parseInt, err := strconv.ParseInt(string(b[i+1:i+4]), 8, 64)
		if err != nil {
			panic(err)
		}
		c[j] = byte(parseInt)
		i += 4
	}

	return c[:j]
}
