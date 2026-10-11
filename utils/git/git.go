package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const maxDiffBytes = 4 << 20

var aheadBehindPattern = regexp.MustCompile(`\[ahead ([0-9]+)(?:, behind ([0-9]+))?\]|\[behind ([0-9]+)(?:, ahead ([0-9]+))?\]`)

type FileStatus struct {
	Path         string `json:"path"`
	OriginalPath string `json:"originalPath,omitempty"`
	Code         string `json:"code"`
	Staged       bool   `json:"staged"`
	Unstaged     bool   `json:"unstaged"`
	Untracked    bool   `json:"untracked"`
	Conflict     bool   `json:"conflict"`
}

type Status struct {
	Root   string       `json:"root"`
	Branch string       `json:"branch"`
	Ahead  int          `json:"ahead"`
	Behind int          `json:"behind"`
	Clean  bool         `json:"clean"`
	Files  []FileStatus `json:"files"`
}

type Diff struct {
	Path     string `json:"path"`
	Original string `json:"original"`
	Working  string `json:"working"`
	Binary   bool   `json:"binary"`
	Deleted  bool   `json:"deleted"`
	ReadOnly bool   `json:"readOnly"`
}

type Worktree struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
}

func ListWorktrees(root string) ([]Worktree, error) {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return nil, err
	}
	output, err := runGit(gitPath, repoRoot, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}

	worktrees := make([]Worktree, 0)
	var current Worktree
	flush := func() {
		if current.Path != "" {
			worktrees = append(worktrees, current)
			current = Worktree{}
		}
	}
	for _, record := range bytes.Split(output, []byte{0}) {
		line := string(record)
		switch {
		case line == "":
			flush()
		case strings.HasPrefix(line, "worktree "):
			flush()
			current.Path = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "branch "):
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		case line == "detached":
			current.Branch = "(detached HEAD)"
		}
	}
	flush()
	sort.Slice(worktrees, func(i, j int) bool {
		return strings.ToLower(worktrees[i].Path) < strings.ToLower(worktrees[j].Path)
	})
	return worktrees, nil
}

func ListBranches(root string) ([]string, error) {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return nil, err
	}
	output, err := runGit(gitPath, repoRoot, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	branches := strings.Fields(string(output))
	sort.Strings(branches)
	return branches, nil
}

func CreateWorktree(root, branch, worktreePath string, createBranch bool) (Worktree, error) {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return Worktree{}, err
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return Worktree{}, errors.New("branch name is empty")
	}
	if _, err := runGit(gitPath, repoRoot, "check-ref-format", "refs/heads/"+branch); err != nil {
		return Worktree{}, fmt.Errorf("invalid branch name %q", branch)
	}
	relPath, err := validateRelativePath(worktreePath)
	if err != nil {
		return Worktree{}, err
	}
	if strings.EqualFold(relPath, ".git") || strings.HasPrefix(strings.ToLower(relPath), ".git"+string(filepath.Separator)) {
		return Worktree{}, errors.New("worktree path cannot be inside .git")
	}
	target := filepath.Join(repoRoot, relPath)
	if err := ensureWorktreeParent(repoRoot, filepath.Dir(relPath)); err != nil {
		return Worktree{}, err
	}
	if _, err := os.Lstat(target); err == nil {
		return Worktree{}, fmt.Errorf("worktree path %q already exists", relPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Worktree{}, err
	}
	if err := excludeWorktreePath(gitPath, repoRoot, relPath); err != nil {
		return Worktree{}, err
	}

	args := []string{"worktree", "add"}
	if createBranch {
		args = append(args, "-b", branch, target)
	} else {
		args = append(args, target, branch)
	}
	if _, err := runGit(gitPath, repoRoot, args...); err != nil {
		return Worktree{}, err
	}
	return Worktree{Path: target, Branch: branch}, nil
}

func excludeWorktreePath(gitPath, repoRoot, relativePath string) error {
	output, err := runGit(gitPath, repoRoot, "rev-parse", "--git-common-dir")
	if err != nil {
		return err
	}
	commonDir := strings.TrimRight(string(output), "\r\n")
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(repoRoot, commonDir)
	}
	excludePath := filepath.Join(commonDir, "info", "exclude")
	existing, err := os.ReadFile(excludePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	rule := "/" + escapeGitExcludePattern(filepath.ToSlash(relativePath)) + "/"
	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSuffix(line, "\r") == rule {
			return nil
		}
	}
	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		existing = append(existing, '\n')
	}
	existing = append(existing, []byte(rule+"\n")...)
	return os.WriteFile(excludePath, existing, 0600)
}

func escapeGitExcludePattern(path string) string {
	return strings.NewReplacer("\\", "\\\\", "*", "\\*", "?", "\\?", "[", "\\[").Replace(path)
}

func ensureWorktreeParent(repoRoot, relativeParent string) error {
	current := repoRoot
	if relativeParent == "." {
		return nil
	}
	for _, segment := range strings.Split(relativeParent, string(filepath.Separator)) {
		current = filepath.Join(current, segment)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(current, 0700); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("worktree parent %q must be a real directory", current)
		}
	}
	return nil
}

func WorkspaceStatus(root string) (Status, error) {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return Status{}, err
	}

	output, err := runGit(gitPath, repoRoot, "status", "--porcelain=v1", "-z", "--branch", "--untracked-files=all")
	if err != nil {
		return Status{}, err
	}

	status, err := parseStatus(output)
	if err != nil {
		return Status{}, err
	}
	status.Root = repoRoot
	status.Clean = len(status.Files) == 0
	return status, nil
}

func ReadDiff(root, filePath string, originalPaths ...string) (Diff, error) {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return Diff{}, err
	}
	relPath, err := validateRelativePath(filePath)
	if err != nil {
		return Diff{}, err
	}

	result := Diff{Path: filepath.ToSlash(relPath)}
	headPath := relPath
	if len(originalPaths) > 0 && strings.TrimSpace(originalPaths[0]) != "" {
		headPath, err = validateRelativePath(originalPaths[0])
		if err != nil {
			return Diff{}, err
		}
	}
	heading := "HEAD:" + filepath.ToSlash(headPath)
	if _, err := runGit(gitPath, repoRoot, "rev-parse", "--verify", "HEAD"); err == nil {
		original, showErr := runGit(gitPath, repoRoot, "show", heading)
		if showErr == nil {
			if len(original) > maxDiffBytes {
				return Diff{}, fmt.Errorf("file exceeds the %d MiB diff limit", maxDiffBytes>>20)
			}
			result.Original = string(original)
		} else if !isMissingGitObject(showErr) {
			return Diff{}, showErr
		}
	}

	workingPath := filepath.Join(repoRoot, relPath)
	fileInfo, statErr := os.Lstat(workingPath)
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return Diff{}, statErr
	}
	result.Deleted = errors.Is(statErr, os.ErrNotExist)
	var working []byte
	if !result.Deleted {
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			linkTarget, linkErr := os.Readlink(workingPath)
			if linkErr != nil {
				return Diff{}, linkErr
			}
			working = []byte(linkTarget)
			result.ReadOnly = true
		} else if !fileInfo.Mode().IsRegular() {
			return Diff{}, fmt.Errorf("cannot diff non-regular file %q", relPath)
		} else {
			working, err = os.ReadFile(workingPath)
			if err != nil {
				return Diff{}, err
			}
		}
		if len(working) > maxDiffBytes {
			return Diff{}, fmt.Errorf("file exceeds the %d MiB diff limit", maxDiffBytes>>20)
		}
		result.Working = string(working)
	}

	result.Binary = bytes.IndexByte([]byte(result.Original), 0) >= 0 || bytes.IndexByte(working, 0) >= 0
	if result.Binary {
		result.Original = ""
		result.Working = ""
		result.ReadOnly = true
	}
	return result, nil
}

func StageFile(root, filePath string) error {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return err
	}
	relPath, err := validateRelativePath(filePath)
	if err != nil {
		return err
	}
	_, err = runGit(gitPath, repoRoot, "add", "--", filepath.ToSlash(relPath))
	return err
}

func UnstageFile(root, filePath string) error {
	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		return err
	}
	relPath, err := validateRelativePath(filePath)
	if err != nil {
		return err
	}
	_, restoreErr := runGit(gitPath, repoRoot, "restore", "--staged", "--", filepath.ToSlash(relPath))
	if restoreErr == nil {
		return nil
	}
	if _, headErr := runGit(gitPath, repoRoot, "rev-parse", "--verify", "HEAD"); headErr != nil {
		_, rmErr := runGit(gitPath, repoRoot, "rm", "--cached", "--", filepath.ToSlash(relPath))
		if rmErr == nil {
			return nil
		}
	}
	return restoreErr
}

func checkedRepository(root string) (gitPath, repoRoot string, err error) {
	gitPath, err = exec.LookPath("git")
	if err != nil {
		return "", "", fmt.Errorf("git executable not found: %w", err)
	}
	if strings.TrimSpace(root) == "" {
		return "", "", errors.New("Git workspace root is empty")
	}

	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	gitEntry := filepath.Join(absoluteRoot, ".git")
	if _, err := os.Stat(gitEntry); err != nil {
		return "", "", fmt.Errorf("project root %q is not a Git working tree", absoluteRoot)
	}

	output, err := runGit(gitPath, absoluteRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}
	repoRoot = filepath.Clean(strings.TrimSpace(string(output)))
	absResolved, absErr := filepath.EvalSymlinks(absoluteRoot)
	repoResolved, repoErr := filepath.EvalSymlinks(repoRoot)
	if absErr == nil && repoErr == nil {
		absoluteRoot = filepath.Clean(absResolved)
		repoRoot = filepath.Clean(repoResolved)
	}
	if !samePath(absoluteRoot, repoRoot) {
		return "", "", fmt.Errorf("project root %q is not the Git root %q", absoluteRoot, repoRoot)
	}
	return gitPath, repoRoot, nil
}

func samePath(left, right string) bool {
	if os.PathSeparator == '\\' {
		left = strings.ToLower(left)
		right = strings.ToLower(right)
	}
	return filepath.Clean(left) == filepath.Clean(right)
}

func validateRelativePath(path string) (string, error) {
	if path == "" || filepath.IsAbs(path) {
		return "", errors.New("file path must be relative to the Git workspace")
	}

	cleanPath := filepath.Clean(filepath.FromSlash(path))
	if cleanPath == "." || cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", errors.New("file path must stay inside the Git workspace")
	}
	return cleanPath, nil
}

func runGit(gitPath, root string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commandArgs := append([]string{"-C", root}, args...)
	cmd := exec.CommandContext(ctx, gitPath, commandArgs...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return output, nil
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("git command timed out")
	}
	message := strings.TrimSpace(string(output))
	if message == "" {
		message = err.Error()
	}
	return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), message)
}

func parseStatus(output []byte) (Status, error) {
	result := Status{Files: make([]FileStatus, 0)}
	entries := bytes.Split(output, []byte{0})
	for index := 0; index < len(entries); index++ {
		entry := entries[index]
		if len(entry) == 0 {
			continue
		}
		if bytes.HasPrefix(entry, []byte("## ")) {
			parseBranchHeader(&result, string(entry[3:]))
			continue
		}
		if len(entry) < 4 || entry[2] != ' ' {
			return Status{}, fmt.Errorf("unrecognized git status record %q", entry)
		}

		indexCode := entry[0]
		worktreeCode := entry[1]
		file := FileStatus{
			Code:      string([]byte{indexCode, worktreeCode}),
			Path:      string(entry[3:]),
			Untracked: indexCode == '?' && worktreeCode == '?',
			Staged:    indexCode != ' ' && indexCode != '?',
			Unstaged:  worktreeCode != ' ' && indexCode != '?',
			Conflict:  indexCode == 'U' || worktreeCode == 'U' || string([]byte{indexCode, worktreeCode}) == "AA" || string([]byte{indexCode, worktreeCode}) == "DD",
		}
		if indexCode == 'R' || indexCode == 'C' || worktreeCode == 'R' || worktreeCode == 'C' {
			if index+1 >= len(entries) {
				return Status{}, fmt.Errorf("missing original path for git status record %q", entry)
			}
			index++
			file.OriginalPath = string(entries[index])
		}
		result.Files = append(result.Files, file)
	}
	sort.SliceStable(result.Files, func(i, j int) bool {
		return result.Files[i].Path < result.Files[j].Path
	})
	return result, nil
}

func parseBranchHeader(result *Status, header string) {
	branch := strings.TrimSpace(strings.SplitN(header, "...", 2)[0])
	branch = strings.TrimPrefix(branch, "No commits yet on ")
	if branch == "HEAD (no branch)" {
		branch = "(detached HEAD)"
	}
	result.Branch = branch

	match := aheadBehindPattern.FindStringSubmatch(header)
	if len(match) == 5 {
		if match[1] != "" {
			result.Ahead, _ = strconv.Atoi(match[1])
			result.Behind, _ = strconv.Atoi(match[2])
		} else {
			result.Behind, _ = strconv.Atoi(match[3])
			result.Ahead, _ = strconv.Atoi(match[4])
		}
	}
}

func isMissingGitObject(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "does not exist in") ||
		strings.Contains(message, "exists on disk, but not in") ||
		strings.Contains(message, "invalid object name") ||
		strings.Contains(message, "needed a single revision") ||
		strings.Contains(message, "path '") && strings.Contains(message, "does not exist")
}
