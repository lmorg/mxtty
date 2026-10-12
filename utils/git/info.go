package git

import (
	"bytes"
	_ "embed"
	"errors"
	"strings"
	"text/template"

	"github.com/lmorg/ttyphoon/utils/markdownutil"
)

//go:embed info.md
var infoTemplate string

type InfoDiffFile struct {
	Filename string
	Diff     string
}

type Info struct {
	RepoWorktreePath string
	RepoStatus       string
	RepoDiffFiles    []InfoDiffFile
	GitExeVersion    string
	GitExePath       string
}

func WorktreeInfo(root, worktreePath string) (Info, error) {
	_, repoRoot, err := checkedRepository(root)
	if err != nil {
		return Info{}, err
	}
	target, err := resolveExistingPathPrefix(worktreePath)
	if err != nil {
		return Info{}, err
	}
	worktrees, err := ListWorktrees(repoRoot)
	if err != nil {
		return Info{}, err
	}
	registered := false
	for _, worktree := range worktrees {
		listedPath, pathErr := resolveExistingPathPrefix(worktree.Path)
		if pathErr == nil && samePath(listedPath, target) {
			registered = true
			break
		}
	}
	if !registered {
		return Info{}, errors.New("selected path is not a registered worktree for the current repository")
	}

	gitPath, worktreeRoot, err := checkedRepository(target)
	if err != nil {
		return Info{}, err
	}
	info := Info{RepoWorktreePath: worktreeRoot, GitExePath: gitPath}

	version, err := runGit(gitPath, worktreeRoot, "--version")
	if err != nil {
		return Info{}, err
	}
	info.GitExeVersion = strings.TrimSpace(string(version))

	status, err := runGit(gitPath, worktreeRoot, "status")
	if err != nil {
		return Info{}, err
	}
	info.RepoStatus = strings.TrimSpace(string(status))

	// Without a HEAD commit there is nothing to compare against, so show staged changes.
	base := []string{"HEAD"}
	if _, err := runGit(gitPath, worktreeRoot, "rev-parse", "--verify", "HEAD"); err != nil {
		base = []string{"--cached"}
	}
	names, err := runGit(gitPath, worktreeRoot, append(append([]string{"diff", "--name-only", "-z"}, base...), "--")...)
	if err != nil {
		return Info{}, err
	}
	for _, name := range bytes.Split(names, []byte{0}) {
		if len(name) == 0 {
			continue
		}
		args := append(append([]string{"diff"}, base...), "--", string(name))
		diff, err := runGit(gitPath, worktreeRoot, args...)
		if err != nil {
			return Info{}, err
		}
		if len(diff) > maxDiffBytes {
			diff = append(diff[:maxDiffBytes], []byte("\n… diff truncated")...)
		}
		info.RepoDiffFiles = append(info.RepoDiffFiles, InfoDiffFile{
			Filename: string(name),
			Diff:     strings.TrimRight(string(diff), "\n"),
		})
	}
	return info, nil
}

func RenderWorktreeInfo(root, worktreePath string) (string, error) {
	info, err := WorktreeInfo(root, worktreePath)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("info").Funcs(markdownutil.TemplateFuncMap()).Parse(infoTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, info); err != nil {
		return "", err
	}
	return buf.String(), nil
}
