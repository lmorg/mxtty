package debug

import (
	"encoding/json"
	"log"
	"runtime"
)

const Enabled = true

func Log(v any) {
	if !Enabled {
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	_, file, line, _ := runtime.Caller(1)

	log.Printf("DEBUG: %s:%d: %s", file, line, string(b))
}
