package types

import (
	"encoding/json"
	"log"
	"strings"
)

type ApcSlice struct {
	slice []string
	kv    map[string]string
}

func NewApcSlice(apc []rune) *ApcSlice {
	s := string(apc)
	as := new(ApcSlice)

	slice := strings.Split(s, ";")
	if len(slice) > 3 {
		as.slice = slice[:2]
		l := len(slice[0]) + len(slice[1]) + 2
		as.slice = append(as.slice, s[l:])
	} else {
		as.slice = slice
	}

	as.kv = make(map[string]string)
	err := json.Unmarshal([]byte(as.Index(2)), &as.kv)
	if err != nil {
		log.Printf("WARNING: cannot decode APC string '%s': %s", s, err.Error())
		//} else {
		//	log.Printf("DEBUG: APC parameters: %s (%v)", as.Index(2), as.kv)
	}

	return as
}

func (as *ApcSlice) Index(i int) string {
	if len(as.slice) <= i {
		return ""
	}
	return as.slice[i]
}

func (as *ApcSlice) Parameter(key string) string {
	return as.kv[key]
}
