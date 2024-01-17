package ptty

import (
	"os"
	"testing"
)

func Test_multiplexBypass(t *testing.T) {
	type args struct {
		p    *PTY
		fifo *os.File
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multiplexBypass(tt.args.p, tt.args.fifo)
		})
	}
}
