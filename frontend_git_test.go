package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitStatusCacheReusesAndRefreshesSnapshots(t *testing.T) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "--quiet"},
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.invalid"},
	} {
		command := exec.Command(gitPath, append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	filePath := filepath.Join(root, "tracked.txt")
	if err := os.WriteFile(filePath, []byte("head\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "--", "tracked.txt"}, {"commit", "--quiet", "-m", "initial"}} {
		command := exec.Command(gitPath, append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}

	app := &WApp{}
	initial, err := app.GetGitStatus(root)
	if err != nil {
		t.Fatalf("GetGitStatus() error = %v", err)
	}
	if !initial.Clean {
		t.Fatalf("initial status = %+v, want clean", initial)
	}

	if err := os.WriteFile(filePath, []byte("changed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cached, err := app.GetGitStatus(root)
	if err != nil || !cached.Clean {
		t.Fatalf("cached status = %+v, %v; want original clean snapshot", cached, err)
	}

	refreshed, err := app.RefreshGitStatus(root)
	if err != nil {
		t.Fatalf("RefreshGitStatus() error = %v", err)
	}
	if refreshed.Clean || len(refreshed.Files) != 1 || !refreshed.Files[0].Unstaged {
		t.Fatalf("refreshed status = %+v, want one unstaged file", refreshed)
	}
	cached, err = app.GetGitStatus(root)
	if err != nil || cached.Clean || len(cached.Files) != 1 {
		t.Fatalf("GetGitStatus() after refresh = %+v, %v", cached, err)
	}
}
