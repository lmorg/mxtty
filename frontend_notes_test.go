package main

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/lmorg/ttyphoon/utils/notes"
)

func TestNotesContextChangesBeforeScanAndCancelsOldScan(t *testing.T) {
	app := &WApp{}
	oldContext := notes.FileContext{PathProjectRoot: t.TempDir(), GroupName: "Old"}
	newContext := notes.FileContext{PathProjectRoot: t.TempDir(), GroupName: "New"}
	app.SyncNotesContext(oldContext)
	t.Cleanup(func() { app.notesIndexUnsubscribe(); app.notesFileIndex.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.notesListCancel = cancel
	sequence := app.notesListSeq

	app.SyncNotesContext(newContext)
	if ctx.Err() != context.Canceled {
		t.Fatal("workspace change did not cancel the previous scan")
	}
	if app.GetCurrentProject() != newContext.PathProjectRoot || app.GetCurrentGroupName() != "New" {
		t.Fatal("workspace context was not updated before scanning")
	}
	if app.notesListSeq <= sequence {
		t.Fatal("workspace change did not invalidate the old scan generation")
	}
	if files, err := app.ListFilesForContext(oldContext); !errors.Is(err, context.Canceled) || files != nil {
		t.Fatalf("stale context listing = %v, %v; want cancellation", files, err)
	}
	sequence = app.notesListSeq
	app.SyncNotesContext(newContext)
	if app.notesListSeq != sequence {
		t.Fatal("unchanged context invalidated the scan generation")
	}
}

func TestNotesFilesSnapshotIsNonblockingAndDoesNotRescan(t *testing.T) {
	started := make(chan struct{}, 1)
	finish := make(chan struct{})
	events := make(chan notes.FileIndexEvent, 4)
	index := notes.NewFileIndex(context.Background(), 4, 1, func(ctx context.Context, root string) notes.RootScanResult {
		started <- struct{}{}
		select {
		case <-finish:
			return notes.RootScanResult{Files: []string{"ready.md"}}
		case <-ctx.Done():
			return notes.RootScanResult{Err: ctx.Err()}
		}
	}, func(event notes.FileIndexEvent) { events <- event })
	t.Cleanup(index.Close)
	app := &WApp{notesFileIndex: index}
	roots := notes.FileContext{PathProjectRoot: t.TempDir(), GroupName: "Workspace"}
	app.SyncNotesContext(roots)
	<-started
	snapshot, err := app.GetNotesFilesSnapshot(roots, false)
	if err != nil || snapshot.State != "scanning" {
		t.Fatalf("initial snapshot = %+v, %v", snapshot, err)
	}
	close(finish)
	select {
	case <-events:
	case <-time.After(3 * time.Second):
		t.Fatal("snapshot was not published")
	}
	snapshot, err = app.GetNotesFilesSnapshot(roots, false)
	if err != nil || !reflect.DeepEqual(snapshot.Files, []string{"$PROJECT/ready.md"}) || snapshot.State != "ready" {
		t.Fatalf("published snapshot = %+v, %v", snapshot, err)
	}
	select {
	case <-started:
		t.Fatal("reading a cached snapshot started another scan")
	default:
	}
}

func TestComposeNotesFilesSnapshotMapsSectionsAndReportsTruncation(t *testing.T) {
	roots := notes.FileContext{PathProjectRoot: t.TempDir(), PathGlobalNotes: t.TempDir(), PathUserNotes: t.TempDir()}
	result := composeNotesFilesSnapshot(roots, []notes.FileIndexSnapshot{
		{Root: roots.PathProjectRoot, Files: []string{"src/main.go"}, State: "truncated", Truncated: true},
		{Root: roots.PathGlobalNotes, Files: []string{"global.md", "group/note.md"}, State: "ready"},
		{Root: roots.PathUserNotes, Files: []string{"todo.md", "nested/note.md"}, State: "ready"},
	})
	want := []string{"$GLOBAL/global.md", "$NOTES/todo.md", "$PROJECT/src/main.go"}
	if !reflect.DeepEqual(result.Files, want) || result.State != "truncated" || !result.Truncated {
		t.Fatalf("snapshot = %+v, want sections %v with truncation", result, want)
	}
}

func TestSaveBinaryFilePublishesWriteThroughHint(t *testing.T) {
	app := &WApp{}
	path := filepath.Join(t.TempDir(), "clipboard.png")
	var changed []string
	unsubscribe := notes.SubscribeFileChanges(func(paths ...string) { changed = append(changed, paths...) })
	defer unsubscribe()
	if err := app.SaveBinaryFile(path, "aW1hZ2U="); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(changed, []string{path}) {
		t.Fatalf("file change hints = %v, want %s", changed, path)
	}
}
