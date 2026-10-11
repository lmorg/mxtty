package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initGitTestRepo(t *testing.T) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	root := t.TempDir()
	runGitTestCommand(t, gitPath, root, "init", "--quiet")
	runGitTestCommand(t, gitPath, root, "config", "user.name", "Test User")
	runGitTestCommand(t, gitPath, root, "config", "user.email", "test@example.invalid")
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("head contents\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "add", "--", "tracked file.txt")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "-m", "initial")
	return root
}

func runGitTestCommand(t *testing.T, gitPath, root string, args ...string) {
	t.Helper()
	commandArgs := append([]string{"-C", root}, args...)
	cmd := exec.Command(gitPath, commandArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func TestWorkspaceStatusAndStageLifecycle(t *testing.T) {
	root := initGitTestRepo(t)
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("working contents\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "new file.txt"), []byte("new contents\n"), 0600); err != nil {
		t.Fatal(err)
	}

	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatalf("WorkspaceStatus() error = %v", err)
	}
	if status.Root != resolveTestPath(t, root) || status.Clean || len(status.Files) != 2 {
		t.Fatalf("WorkspaceStatus() = %+v, want root, dirty state, and two files", status)
	}
	files := make(map[string]FileStatus, len(status.Files))
	for _, file := range status.Files {
		files[file.Path] = file
	}
	if !files["tracked file.txt"].Unstaged || files["tracked file.txt"].Untracked {
		t.Fatalf("tracked file status = %+v, want unstaged tracked change", files["tracked file.txt"])
	}
	if !files["new file.txt"].Untracked || files["new file.txt"].Staged {
		t.Fatalf("new file status = %+v, want untracked change", files["new file.txt"])
	}

	diff, err := ReadDiff(root, "tracked file.txt")
	if err != nil {
		t.Fatalf("ReadDiff(tracked) error = %v", err)
	}
	if diff.Original != "head contents\n" || diff.Working != "working contents\n" {
		t.Fatalf("ReadDiff(tracked) = %+v", diff)
	}

	diff, err = ReadDiff(root, "new file.txt")
	if err != nil {
		t.Fatalf("ReadDiff(untracked) error = %v", err)
	}
	if diff.Original != "" || diff.Working != "new contents\n" {
		t.Fatalf("ReadDiff(untracked) = %+v, want empty original and working contents", diff)
	}

	if err := StageFile(root, "new file.txt"); err != nil {
		t.Fatalf("StageFile() error = %v", err)
	}
	status, err = WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Files) != 2 || !status.Files[0].Staged && !status.Files[1].Staged {
		t.Fatalf("staged status = %+v, want staged file", status.Files)
	}

	if err := UnstageFile(root, "new file.txt"); err != nil {
		t.Fatalf("UnstageFile() error = %v", err)
	}
	status, err = WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Files) != 2 || !filesForPath(status.Files, "new file.txt").Untracked {
		t.Fatalf("unstaged status = %+v, want new file restored to untracked", status.Files)
	}
}

func TestWorkspaceStatusAcceptsWorktreeRoot(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	worktree := filepath.Join(t.TempDir(), "linked-worktree")
	runGitTestCommand(t, gitPath, root, "worktree", "add", "--quiet", "-b", "linked", worktree)

	status, err := WorkspaceStatus(worktree)
	if err != nil {
		t.Fatalf("WorkspaceStatus(worktree) error = %v", err)
	}
	if status.Root != resolveTestPath(t, worktree) || status.Branch != "linked" {
		t.Fatalf("WorkspaceStatus(worktree) = %+v, want worktree root and linked branch", status)
	}
	info, err := os.Stat(filepath.Join(worktree, ".git"))
	if err != nil || info.IsDir() {
		t.Fatalf("worktree .git entry = %v, %v; want a gitdir file", info, err)
	}
}

func TestListBranchesAndWorktreesFromEitherRootAndCreateWorktrees(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	linkedRoot := filepath.Join(t.TempDir(), " linked worktree")
	runGitTestCommand(t, gitPath, root, "worktree", "add", "--quiet", "-b", "zeta", linkedRoot)
	runGitTestCommand(t, gitPath, root, "branch", "alpha")

	for _, queryRoot := range []string{root, linkedRoot} {
		branches, err := ListBranches(queryRoot)
		if err != nil {
			t.Fatalf("ListBranches(%q) error = %v", queryRoot, err)
		}
		if len(branches) != 3 || branches[0] != "alpha" || branches[1] != "main" || branches[2] != "zeta" {
			t.Fatalf("ListBranches(%q) = %v, want sorted local branches", queryRoot, branches)
		}
		worktrees, err := ListWorktrees(queryRoot)
		if err != nil {
			t.Fatalf("ListWorktrees(%q) error = %v", queryRoot, err)
		}
		if len(worktrees) != 2 {
			t.Fatalf("ListWorktrees(%q) = %+v, want both worktrees", queryRoot, worktrees)
		}
		if !containsWorktreePath(worktrees, resolveTestPath(t, linkedRoot)) {
			t.Fatalf("ListWorktrees(%q) = %+v, missing exact whitespace-containing path %q", queryRoot, worktrees, linkedRoot)
		}
	}

	createdPath := filepath.Join(resolveTestPath(t, root), ".git_worktrees", "alpha")
	worktree, err := CreateWorktree(root, "alpha", filepath.Join(".git_worktrees", "alpha"), false)
	if err != nil {
		t.Fatalf("CreateWorktree(existing branch) error = %v", err)
	}
	if worktree.Path != createdPath || worktree.Branch != "alpha" {
		t.Fatalf("CreateWorktree(existing branch) = %+v", worktree)
	}
	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatalf("WorkspaceStatus() after nested worktree creation error = %v", err)
	}
	if !status.Clean {
		t.Fatalf("WorkspaceStatus() after nested worktree creation = %+v, want clean parent", status)
	}

	newBranchPath := filepath.Join(".git_worktrees", "feature", "new")
	worktree, err = CreateWorktree(linkedRoot, "feature/new", newBranchPath, true)
	if err != nil {
		t.Fatalf("CreateWorktree(new branch from linked root) error = %v", err)
	}
	if worktree.Branch != "feature/new" {
		t.Fatalf("CreateWorktree(new branch) = %+v", worktree)
	}
}

func containsWorktreePath(worktrees []Worktree, path string) bool {
	for _, worktree := range worktrees {
		if worktree.Path == path {
			return true
		}
	}
	return false
}

func TestCreateWorktreeRejectsUnsafePathsAndBranches(t *testing.T) {
	root := initGitTestRepo(t)
	if _, err := CreateWorktree(root, "main", "../outside", false); err == nil {
		t.Fatal("CreateWorktree() error = nil, want path traversal rejection")
	}
	if _, err := CreateWorktree(root, "invalid branch", "safe/path", true); err == nil {
		t.Fatal("CreateWorktree() error = nil, want invalid branch rejection")
	}
}

func TestStageFileRejectsPathsOutsideWorkspace(t *testing.T) {
	root := initGitTestRepo(t)
	if err := StageFile(root, "../outside.txt"); err == nil {
		t.Fatal("StageFile() error = nil, want traversal rejection")
	}
	if _, err := ReadDiff(root, filepath.Join(root, "tracked file.txt")); err == nil {
		t.Fatal("ReadDiff() error = nil, want absolute path rejection")
	}
}

func TestReadDiffTreatsSymlinkAsReadOnly(t *testing.T) {
	root := initGitTestRepo(t)
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape.txt")); err != nil {
		t.Skipf("cannot create symlinks: %v", err)
	}
	diff, err := ReadDiff(root, "escape.txt")
	if err != nil {
		t.Fatalf("ReadDiff(symlink) error = %v", err)
	}
	if !diff.ReadOnly || diff.Working != outside {
		t.Fatalf("ReadDiff(symlink) = %+v, want read-only link target text", diff)
	}
}

func TestUnstageFileOnUnbornBranch(t *testing.T) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	root := t.TempDir()
	runGitTestCommand(t, gitPath, root, "init", "--quiet")
	if err := os.WriteFile(filepath.Join(root, "first.txt"), []byte("first"), 0600); err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "add", "--", "first.txt")

	if err := UnstageFile(root, "first.txt"); err != nil {
		t.Fatalf("UnstageFile() on unborn branch error = %v", err)
	}
	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Files) != 1 || !status.Files[0].Untracked {
		t.Fatalf("status after unstage = %+v, want first.txt untracked", status.Files)
	}
}

func TestWorkspacePreservesWhitespaceInFilePaths(t *testing.T) {
	root := initGitTestRepo(t)
	const filePath = " leading name.txt "
	if err := os.WriteFile(filepath.Join(root, filePath), []byte("contents"), 0600); err != nil {
		t.Fatal(err)
	}
	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if !filesForPath(status.Files, filePath).Untracked {
		t.Fatalf("status = %+v, want exact whitespace path %q", status.Files, filePath)
	}
	if err := StageFile(root, filePath); err != nil {
		t.Fatalf("StageFile(%q) error = %v", filePath, err)
	}
}

func TestReadDiffUsesOriginalPathForStagedRename(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	runGitTestCommand(t, gitPath, root, "mv", "tracked file.txt", "renamed file.txt")

	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Files) != 1 || status.Files[0].Path != "renamed file.txt" || status.Files[0].OriginalPath != "tracked file.txt" {
		t.Fatalf("rename status = %+v", status.Files)
	}
	diff, err := ReadDiff(root, status.Files[0].Path, status.Files[0].OriginalPath)
	if err != nil {
		t.Fatalf("ReadDiff(rename) error = %v", err)
	}
	if diff.Original != "head contents\n" || diff.Working != "head contents\n" {
		t.Fatalf("ReadDiff(rename) = %+v, want content from original HEAD path", diff)
	}
}

func filesForPath(files []FileStatus, path string) FileStatus {
	for _, file := range files {
		if file.Path == path {
			return file
		}
	}
	return FileStatus{}
}

func resolveTestPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(resolved)
}
