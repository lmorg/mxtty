# 0062 - Notes Git panel is scoped to the current project root

## Status

Accepted and implemented on 2026-10-09.

## Context

Notes needs Git status and file-level actions without making Git command
execution or repository parsing part of the frontend. The active Notes project
root may be a normal repository root, a linked worktree root, or a non-Git
`project.ttyphoon` root. A Git panel should not appear for the latter or when
Git is unavailable.

## Decision

- Put Git command execution and status/diff parsing in `utils/git`; expose thin
  Wails wrappers from `frontend.go`.
- Enable the Tools Git tab only when the current project root is exactly a Git
  working-tree root and the host has a `git` executable. Accept both `.git`
  directories and worktree `.git` files. The workspace dropdown contains only
  the current project root/worktree; do not enumerate other repositories.
- Render branch state, ahead/behind counts, and changed paths in the Git pane.
  Selecting a path opens a side-by-side Monaco diff comparing `HEAD` with the
  working tree in a dedicated top-level Notes `Diff` tab. The modified side is
  editable for regular text files; binary and symbolic-link diffs are read-only.
- Open the selected path through Notes' normal document loader. The diff's
  modified model stays synchronized with the regular Notes editor model, so
  filename display, recent documents, history, autosave, and save hooks remain
  shared with ordinary files.
- Keep separate Monaco editor and diff-editor factories because Monaco exposes
  different widget APIs, but share keyboard/completion, clipboard, and context
  menu hooks between them.
- Support file-level stage and unstage actions, but no commit action. Save dirty
  diff edits before staging or unstaging the same file. Select operations use
  argument-separated Git invocations and paths relative to the validated root.
- Handle filenames as NUL-delimited Git status records, including spaces and
  rename source paths. Cap text diff inputs and reject path traversal and
  symlink escapes.

## Consequences

- The panel follows the active project context. A root change refreshes status;
  a non-Git root hides the tab.
- Selecting a Git file makes it the active Notes document. The normal Notes
  filename, recent list, and history apply; the diff overlays the regular editor
  while keeping its model synchronized. Leaving the diff reveals the regular
  editor for the same file. The Diff tab is present only while a Git diff is
  active. Cmd/Ctrl+S keeps the normal format-and-save behavior.
- If Monaco is disabled, opening a Git file still loads it into Notes and
  updates history/recents, but the side-by-side diff is unavailable.
- Git operations are currently read, stage, and unstage only. Commit, branch
  management, multi-repository discovery, and persistent workspace lists are
  out of scope.

## Reference

- `utils/git/git.go` - repository validation, status, diff, stage, and unstage
- `frontend.go` - Wails wrappers
- `frontend/src/notes-git-panel.js` - status panel and file actions
- `frontend/src/monaco_adapter.js` - shared editor hooks and diff-editor facade
- `frontend/src/notes.js`, `frontend/src/notes.css` - tab, workspace, and editor
  integration
- `utils/git/git_test.go`, `frontend/src/notes-git-panel.test.js` - behavior
  coverage
- ADR 0013 - Wails binding checklist
- ADR 0018 - Monaco provider architecture
