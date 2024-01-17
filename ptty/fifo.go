package ptty

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/lmorg/mxtty/app"
)

func (p *PTY) createFifo() string {
	filename := fmt.Sprintf("%s/%s-%d.fifo", os.TempDir(), app.Name, os.Getpid())

	err := syscall.Mkfifo(filename, 0640)
	if err != nil {
		panic(err)
	}

	log.Printf("DEBUG: FIFO created at %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	go p.read(f)

	return filename
}
