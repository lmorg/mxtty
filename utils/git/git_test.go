package git

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestBranchCheckoutAndFetch(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckoutBranch(root, "topic/new", true); err != nil {
		t.Fatal(err)
	}
	status, err := WorkspaceStatus(root)
	if err != nil || status.Branch != "topic/new" {
		t.Fatalf("created branch = %+v, %v", status, err)
	}
	remote := initGitTestRepo(t)
	runGitTestCommand(t, gitPath, remote, "checkout", "--quiet", "-b", "topic/new")
	runGitTestCommand(t, gitPath, remote, "commit", "--quiet", "--allow-empty", "-m", "remote update")
	runGitTestCommand(t, gitPath, root, "remote", "add", "origin", remote)
	head, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	result, err := FetchBranch(root, "topic/new")
	if err != nil || result["error"] != "" {
		t.Fatalf("fetch = %v, %v", result, err)
	}
	after, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil || string(head) != string(after) {
		t.Fatalf("fetch changed HEAD: %s, %v", after, err)
	}
	remoteHead, err := exec.Command(gitPath, "-C", remote, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := exec.Command(gitPath, "-C", root, "rev-parse", "refs/remotes/origin/topic/new").Output()
	if err != nil || string(fetched) != string(remoteHead) {
		t.Fatalf("fetched tip = %s, %v", fetched, err)
	}
	runGitTestCommand(t, gitPath, remote, "checkout", "--quiet", "-b", "remote/topic")
	runGitTestCommand(t, gitPath, remote, "commit", "--quiet", "--allow-empty", "-m", "upstream update")
	runGitTestCommand(t, gitPath, root, "config", "branch.topic/new.remote", "origin")
	runGitTestCommand(t, gitPath, root, "config", "branch.topic/new.merge", "refs/heads/remote/topic")
	result, err = FetchBranch(root, "topic/new")
	if err != nil || result["error"] != "" {
		t.Fatalf("upstream fetch = %v, %v", result, err)
	}
	remoteHead, err = exec.Command(gitPath, "-C", remote, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	fetched, err = exec.Command(gitPath, "-C", root, "rev-parse", "refs/remotes/origin/remote/topic").Output()
	if err != nil || string(fetched) != string(remoteHead) {
		t.Fatalf("upstream tip = %s, %v", fetched, err)
	}
	if err := CheckoutBranch(root, "other", true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("other branch\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "add", "tracked file.txt")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "-m", "other content")
	if err := CheckoutBranch(root, "topic/new", false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("unsaved edits\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := CheckoutBranch(root, "other", false); err == nil {
		t.Fatal("checkout overwrote dirty content")
	}
	status, err = WorkspaceStatus(root)
	if err != nil || status.Branch != "topic/new" {
		t.Fatalf("failed checkout changed branch: %+v, %v", status, err)
	}
	runGitTestCommand(t, gitPath, root, "config", "branch.topic/new.merge", "refs/heads/missing")
	result, err = FetchBranch(root, "topic/new")
	if err != nil || result["error"] == "" || result["stderr"] == "" {
		t.Fatalf("fetch failure output = %v, %v", result, err)
	}
	for _, invalid := range []string{"invalid branch", "--detach", "HEAD~1"} {
		if err := CheckoutBranch(root, invalid, false); err == nil {
			t.Fatalf("accepted %q", invalid)
		}
	}
}

func TestListCommitsIncludesMetadataAndGraph(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "commit", "--amend", "--quiet", "-m", "initial", "-m", "Body-only search phrase")
	commits, err := ListCommits(root, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 || commits[0].Subject != "initial" || commits[0].Author != "Test User" {
		t.Fatalf("ListCommits() = %+v", commits)
	}
	if len(commits[0].Hash) < 7 || len(commits[0].Date) != 19 || !strings.Contains(commits[0].Graph[0], "*") {
		t.Fatalf("commit metadata/graph = %+v", commits[0])
	}
	fullHash, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil || commits[0].FullHash != strings.TrimSpace(string(fullHash)) {
		t.Fatalf("full commit hash = %q, %v", commits[0].FullHash, err)
	}
	if commits[0].Message != "initial\n\nBody-only search phrase" {
		t.Fatalf("commit message = %q", commits[0].Message)
	}
	if _, err := ListCommits(root, 0); err == nil {
		t.Fatal("ListCommits accepted an invalid limit")
	}
}

func TestListCommitsPreservesMergeGraphAndLimit(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	output, err := exec.Command(gitPath, "-C", root, "branch", "--show-current").Output()
	if err != nil {
		t.Fatal(err)
	}
	branch := strings.TrimSpace(string(output))
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", "-b", "feature")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "--allow-empty", "-m", "Feature commit")
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", branch)
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "--allow-empty", "-m", "Main commit")
	runGitTestCommand(t, gitPath, root, "merge", "--quiet", "--no-ff", "feature", "-m", "Merge feature")
	commits, err := ListCommits(root, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 4 || commits[0].Subject != "Merge feature" {
		t.Fatalf("merge history = %+v", commits)
	}
	var graph strings.Builder
	for _, commit := range commits {
		graph.WriteString(strings.Join(commit.Graph, "\n"))
	}
	if !strings.Contains(graph.String(), "\\") || !strings.Contains(graph.String(), "/") {
		t.Fatalf("merge connectors missing from graph: %q", graph.String())
	}
	limited, err := ListCommits(root, 2)
	if err != nil || len(limited) != 2 {
		t.Fatalf("limited history = %+v, %v", limited, err)
	}
}

func TestListCommitsScopesHistoryToHead(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	output, err := exec.Command(gitPath, "-C", root, "branch", "--show-current").Output()
	if err != nil {
		t.Fatal(err)
	}
	branch := strings.TrimSpace(string(output))
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", "-b", "unmerged")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "--allow-empty", "-m", "Unmerged commit")
	runGitTestCommand(t, gitPath, root, "tag", "unmerged-tag")
	runGitTestCommand(t, gitPath, root, "update-ref", "refs/remotes/origin/unmerged", "HEAD")
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", branch)

	commits, err := ListCommits(root, 100)
	if err != nil || len(commits) != 1 || commits[0].Subject != "initial" {
		t.Fatalf("current branch history = %+v, %v; want only initial commit", commits, err)
	}
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", "--detach", "unmerged")
	commits, err = ListCommits(root, 100)
	if err != nil || len(commits) != 2 || commits[0].Subject != "Unmerged commit" {
		t.Fatalf("detached HEAD history = %+v, %v; want detached commit's ancestry", commits, err)
	}
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", "--orphan", "unborn")
	commits, err = ListCommits(root, 100)
	if err != nil || len(commits) != 0 {
		t.Fatalf("unborn branch history = %+v, %v; want no commits", commits, err)
	}
}

func TestWorkspaceStatusDoesNotRewriteIndex(t *testing.T) {
	root := initGitTestRepo(t)
	t.Setenv("GIT_OPTIONAL_LOCKS", "1")
	indexPath := filepath.Join(root, ".git", "index")
	before, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	trackedPath := filepath.Join(root, "tracked file.txt")
	stale := time.Now().Add(-time.Hour)
	if err := os.Chtimes(trackedPath, stale, stale); err != nil {
		t.Fatal(err)
	}
	status, err := WorkspaceStatus(root)
	if err != nil || !status.Clean {
		t.Fatalf("timestamp-only status = %+v, %v; want clean", status, err)
	}
	if err := os.WriteFile(trackedPath, []byte("changed contents\n"), 0600); err != nil {
		t.Fatal(err)
	}
	status, err = WorkspaceStatus(root)
	if err != nil || status.Clean || !filesForPath(status.Files, "tracked file.txt").Unstaged {
		t.Fatalf("changed-content status = %+v, %v; want an unstaged change", status, err)
	}
	after, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) || !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatal("status refreshed the index on disk; this can retrigger the metadata watch")
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

func TestMetadataDirectoriesIncludeLinkedWorktreeCommonDirectory(t *testing.T) {
	root := initGitTestRepo(t)
	worktree := filepath.Join(t.TempDir(), "linked-worktree")
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	runGitTestCommand(t, gitPath, root, "worktree", "add", "--quiet", "--detach", worktree, "HEAD")

	directories, err := MetadataDirectories(worktree)
	if err != nil {
		t.Fatalf("MetadataDirectories() error = %v", err)
	}
	if len(directories) != 2 {
		t.Fatalf("MetadataDirectories() = %v, want per-worktree and common directories", directories)
	}
	for _, directory := range directories {
		if directory == filepath.Join(worktree, ".git") {
			t.Fatalf("MetadataDirectories() returned the .git file path: %v", directories)
		}
		if _, err := os.Stat(directory); err != nil {
			t.Fatalf("metadata directory %q is not accessible: %v", directory, err)
		}
	}
}

func TestBrowserLinkForRemote(t *testing.T) {
	tests := []struct {
		name     string
		remote   string
		branch   string
		provider string
		url      string
	}{
		{"GitHub HTTPS", "https://github.com/org/repo.git", "feature/view", "github", "https://github.com/org/repo/tree/feature/view"},
		{"GitHub SSH", "git@github.com:org/repo.git", "main", "github", "https://github.com/org/repo/tree/main"},
		{"GitLab nested group", "git@gitlab.com:org/group/repo.git", "feature/view", "gitlab", "https://gitlab.com/org/group/repo/-/tree/feature/view"},
		{"Bitbucket", "git@bitbucket.org:org/repo.git", "main", "bitbucket", "https://bitbucket.org/org/repo/src/main/"},
		{"Azure HTTPS", "https://dev.azure.com/org/project/_git/repo", "feature/view", "azuredevops", "https://dev.azure.com/org/project/_git/repo?version=GBfeature%2Fview"},
		{"Azure SSH", "git@ssh.dev.azure.com:v3/org/project/repo", "main", "azuredevops", "https://dev.azure.com/org/project/_git/repo?version=GBmain"},
		{"Azure legacy host", "https://org.visualstudio.com/project/_git/repo", "main", "azuredevops", "https://dev.azure.com/org/project/_git/repo?version=GBmain"},
		{"unsupported host", "git@example.com:org/repo.git", "main", "", ""},
		{"local path", "../repo", "main", "", ""},
		{"detached head", "https://github.com/org/repo.git", "HEAD", "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			link := browserLinkForRemote(test.remote, test.branch)
			if link.Provider != test.provider || link.URL != test.url {
				t.Fatalf("browserLinkForRemote(%q, %q) = %+v, want provider %q and URL %q", test.remote, test.branch, link, test.provider, test.url)
			}
		})
	}
}

func TestRemoteBrowserLinkReadsOriginFromRepository(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	link, err := RemoteBrowserLink(root, "main")
	if err != nil || link.URL != "" {
		t.Fatalf("RemoteBrowserLink() without origin = %+v, %v; want no link", link, err)
	}
	runGitTestCommand(t, gitPath, root, "remote", "add", "origin", "git@github.com:org/repo.git")

	link, err = RemoteBrowserLink(root, "feature/ui")
	if err != nil {
		t.Fatalf("RemoteBrowserLink() error = %v", err)
	}
	if link.Provider != "github" || link.URL != "https://github.com/org/repo/tree/feature/ui" {
		t.Fatalf("RemoteBrowserLink() = %+v, want GitHub feature branch URL", link)
	}
}

func TestBranchDiffFilesAndReadDiffDoNotSwitchBranches(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	branchCommand := exec.Command(gitPath, "-C", root, "branch", "--show-current")
	branchOutput, err := branchCommand.Output()
	if err != nil {
		t.Fatal(err)
	}
	currentBranch := strings.TrimSpace(string(branchOutput))
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", "-b", "compare-base")
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("compare contents\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "branch-only.txt"), []byte("branch file\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "add", "--", "tracked file.txt", "branch-only.txt")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "-m", "comparison base")
	runGitTestCommand(t, gitPath, root, "checkout", "--quiet", currentBranch)
	if err := os.WriteFile(filepath.Join(root, "untracked.txt"), []byte("local file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	files, err := BranchDiffFiles(root, "compare-base")
	if err != nil {
		t.Fatalf("BranchDiffFiles() error = %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("BranchDiffFiles() = %+v, want tracked, branch-only, and untracked paths", files)
	}
	for _, file := range files {
		if !file.CompareOnly {
			t.Errorf("file %q CompareOnly = false", file.Path)
		}
	}
	diff, err := ReadDiff(root, "tracked file.txt", "", "compare-base")
	if err != nil {
		t.Fatalf("ReadDiff(comparison branch) error = %v", err)
	}
	if diff.Original != "compare contents\n" || diff.Working != "head contents\n" {
		t.Fatalf("ReadDiff(comparison branch) = %+v, want comparison base and current working contents", diff)
	}
	status, err := WorkspaceStatus(root)
	if err != nil {
		t.Fatal(err)
	}
	if status.Branch != currentBranch {
		t.Fatalf("WorkspaceStatus().Branch = %q, want unchanged branch %q", status.Branch, currentBranch)
	}
}

func TestCommitComparisonUsesCommitWithoutChangingHead(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	original, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	commit := strings.TrimSpace(string(original))
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("later contents\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGitTestCommand(t, gitPath, root, "add", "--", "tracked file.txt")
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "-m", "later")
	head, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{commit, commit[:7]} {
		files, err := BranchDiffFiles(root, target)
		if err != nil || len(files) != 1 {
			t.Fatalf("commit comparison files = %+v, %v", files, err)
		}
		diff, err := ReadDiff(root, "tracked file.txt", "", target)
		if err != nil || diff.Original != "head contents\n" || diff.Working != "later contents\n" {
			t.Fatalf("commit diff = %+v, %v", diff, err)
		}
	}
	after, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil || string(after) != string(head) {
		t.Fatalf("HEAD changed after comparison: %s, %v", after, err)
	}
	runGitTestCommand(t, gitPath, root, "branch", commit[:8], "HEAD")
	files, err := BranchDiffFiles(root, commit[:8])
	if err != nil || len(files) != 0 {
		t.Fatalf("hex-named branch precedence = %+v, %v", files, err)
	}
	for _, target := range []string{"HEAD~1", "--all", "0000000"} {
		if _, err := BranchDiffFiles(root, target); err == nil {
			t.Fatalf("accepted invalid comparison target %q", target)
		}
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

func TestCreateWorktreeAtCommitCreatesDetachedHistory(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	output, err := exec.Command(gitPath, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	hash := strings.TrimSpace(string(output))
	shortOutput, err := exec.Command(gitPath, "-C", root, "rev-parse", "--short", hash).Output()
	if err != nil {
		t.Fatal(err)
	}
	shortHash := strings.TrimSpace(string(shortOutput))
	runGitTestCommand(t, gitPath, root, "commit", "--quiet", "--allow-empty", "-m", "later")
	target := filepath.Join(t.TempDir(), shortHash)
	worktree, err := CreateWorktree(root, hash, target, false)
	if err != nil {
		t.Fatal(err)
	}
	if worktree.Branch != shortHash {
		t.Fatalf("worktree = %+v", worktree)
	}
	for _, listingRoot := range []string{root, target} {
		listed, err := ListWorktrees(listingRoot)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, entry := range listed {
			if entry.Path == worktree.Path {
				found = true
				if entry.Branch != shortHash {
					t.Fatalf("detached menu label = %q, want %q", entry.Branch, shortHash)
				}
			}
		}
		if !found {
			t.Fatal("detached worktree missing from listing")
		}
	}
	actual, err := exec.Command(gitPath, "-C", target, "rev-parse", "HEAD").Output()
	if err != nil || strings.TrimSpace(string(actual)) != hash {
		t.Fatalf("worktree HEAD = %s, %v", actual, err)
	}
	if err := exec.Command(gitPath, "-C", target, "symbolic-ref", "--quiet", "HEAD").Run(); err == nil {
		t.Fatal("commit worktree is not detached")
	}
	if _, err := CreateWorktree(root, strings.Repeat("0", 40), filepath.Join(t.TempDir(), "invalid"), false); err == nil {
		t.Fatal("accepted a nonexistent commit")
	}
	runGitTestCommand(t, gitPath, root, "branch", "abcdef1")
	hexTarget := filepath.Join(t.TempDir(), "hex-branch")
	if _, err := CreateWorktree(root, "abcdef1", hexTarget, false); err != nil {
		t.Fatal(err)
	}
	branch, err := exec.Command(gitPath, "-C", hexTarget, "symbolic-ref", "--short", "HEAD").Output()
	if err != nil || strings.TrimSpace(string(branch)) != "abcdef1" {
		t.Fatalf("hex branch worktree = %s, %v", branch, err)
	}
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

func TestDefaultWorktreePathUsesOriginOrProjectRootHash(t *testing.T) {
	root := initGitTestRepo(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	rootHash := sha256.Sum256([]byte(resolveTestPath(t, root)))
	path, err := DefaultWorktreePath(root)
	if err != nil {
		t.Fatalf("DefaultWorktreePath() without origin error = %v", err)
	}
	want := filepath.Join(home, "Documents", "ttyphoon", ".git_worktrees", hex.EncodeToString(rootHash[:]))
	if path != want {
		t.Fatalf("DefaultWorktreePath() without origin = %q, want %q", path, want)
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	origin := "https://example.invalid/team/project.git"
	runGitTestCommand(t, gitPath, root, "remote", "add", "origin", origin)
	originHash := sha256.Sum256([]byte(origin))
	path, err = DefaultWorktreePath(root)
	if err != nil {
		t.Fatalf("DefaultWorktreePath() with origin error = %v", err)
	}
	want = filepath.Join(home, "Documents", "ttyphoon", ".git_worktrees", hex.EncodeToString(originHash[:]))
	if path != want {
		t.Fatalf("DefaultWorktreePath() with origin = %q, want %q", path, want)
	}
}

func TestCreateWorktreeAtAbsolutePathOutsideRepository(t *testing.T) {
	root := initGitTestRepo(t)
	worktreePath := filepath.Join(t.TempDir(), "external-worktree")
	worktree, err := CreateWorktree(root, "external", worktreePath, true)
	if err != nil {
		t.Fatalf("CreateWorktree(absolute path) error = %v", err)
	}
	if worktree.Path != resolveTestPath(t, worktreePath) || worktree.Branch != "external" {
		t.Fatalf("CreateWorktree(absolute path) = %+v, want path %q and branch external", worktree, worktreePath)
	}
}

func TestRemoveWorktreePreservesBranch(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	worktreePath := filepath.Join(t.TempDir(), "removable-worktree")
	runGitTestCommand(t, gitPath, root, "worktree", "add", "--quiet", "-b", "removable", worktreePath)
	if err := RemoveWorktree(root, worktreePath); err != nil {
		t.Fatalf("RemoveWorktree() error = %v", err)
	}
	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("worktree path stat error = %v, want not-exist", err)
	}
	output, err := exec.Command(gitPath, "-C", root, "branch", "--list", "removable").CombinedOutput()
	if err != nil {
		t.Fatalf("git branch --list error = %v: %s", err, output)
	}
	if strings.TrimSpace(string(output)) != "removable" {
		t.Fatalf("branch list after RemoveWorktree = %q, want retained branch removable", output)
	}
}

func TestRenderWorktreeInfoKeepsMultilineDiffInOneTableRow(t *testing.T) {
	root := initGitTestRepo(t)
	if err := os.WriteFile(filepath.Join(root, "tracked file.txt"), []byte("head contents\nnew | <line>\n"), 0600); err != nil {
		t.Fatal(err)
	}
	markdown, err := RenderWorktreeInfo(root, root)
	if err != nil {
		t.Fatalf("RenderWorktreeInfo() error = %v", err)
	}
	if !strings.Contains(markdown, "| tracked file\\.txt | <pre><code>diff --git") {
		t.Fatalf("RenderWorktreeInfo() missing diff table row:\n%s", markdown)
	}
	if !strings.Contains(markdown, "&#10;+new &#124; &lt;line&gt;</code></pre> |") {
		t.Fatalf("RenderWorktreeInfo() did not encode multiline diff:\n%s", markdown)
	}
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "| tracked file\\.txt |") && !strings.HasSuffix(line, "</code></pre> |") {
			t.Fatalf("diff row was split across lines: %q", line)
		}
	}
	if !strings.Contains(markdown, "git version") {
		t.Fatalf("RenderWorktreeInfo() missing Git version:\n%s", markdown)
	}
}

func TestRemoveWorktreeRejectsDirtyWorkingTreeWithoutForce(t *testing.T) {
	root := initGitTestRepo(t)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed")
	}
	worktreePath := filepath.Join(t.TempDir(), "dirty-worktree")
	runGitTestCommand(t, gitPath, root, "worktree", "add", "--quiet", "-b", "dirty", worktreePath)
	if err := os.WriteFile(filepath.Join(worktreePath, "tracked file.txt"), []byte("modified\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := RemoveWorktree(root, worktreePath); err == nil {
		t.Fatal("RemoveWorktree() error = nil, want dirty-worktree rejection")
	}
	if _, err := os.Stat(filepath.Join(worktreePath, "tracked file.txt")); err != nil {
		t.Fatalf("dirty worktree should remain after rejected removal: %v", err)
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
