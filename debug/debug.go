package debug

import (
	"encoding/json"
	"log"
	"runtime"
	"strings"

	_ "github.com/lmorg/mxtty/debug/pprof"
)

func Log(v any) {
	if !Enabled {
		return
	}

	switch t := v.(type) {
	case byte:
		v = string(t)

	case []byte:
		v = string(t)

	case []rune:
		if len(t) > 0 && t[0] < 32 {
			break
		}
		v = string(t)
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	pc, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Printf("DEBUG: %s:%d: %s", file, line, string(b))
		return
	}

	fn := runtime.FuncForPC(pc)
	fnName := strings.Replace(fn.Name(), "github.com/lmorg/mxtty/", "", 1)

	pc, _, _, ok = runtime.Caller(2)
	if !ok {
		log.Printf("DEBUG: %s(): %s", fnName, string(b))
		return
	}

	fn = runtime.FuncForPC(pc)
	prevName := strings.Replace(fn.Name(), "github.com/lmorg/mxtty/", "", 1)
	log.Printf("DEBUG: %s() <- %s(): %s", fnName, prevName, string(b))
}
