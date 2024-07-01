package exit

import (
	"os"

	"github.com/lmorg/mxtty/debug/pprof"
)

func Exit(code int) {
	pprof.CleanUp()
	os.Exit(code)
}
