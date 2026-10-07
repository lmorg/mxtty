package dirwatch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchReportsCreationAndStopsOnCancellation(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ready := make(chan error, 1)
	changes := make(chan Change, 128)
	done := make(chan struct{})
	go func() {
		Run(ctx, root, func(string) bool { return false }, func(change Change) { changes <- change }, func(err error) { ready <- err })
		close(done)
	}()
	if err := <-ready; err != nil {
		t.Fatalf("native watcher failed: %v", err)
	}
	file := filepath.Join(root, "created.md")
	if err := os.WriteFile(file, []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case change := <-changes:
			if change.Path != file && !change.Rescan {
				continue
			}
			cancel()
			select {
			case <-done:
				return
			case <-timeout.C:
				t.Fatal("watcher did not stop")
			}
		case <-timeout.C:
			t.Fatal("file creation was not reported")
		}
	}
}
