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
- Cache Git status in Go by canonical project root. Reuse the cached snapshot on
  reads and explicitly refresh after stage/unstage or a file-index watcher
  completion; publish refreshed snapshots to the frontend. Keep a frontend cache
  keyed by project root so switching Tools tabs or worktrees can rehydrate
  immediately.
- Enable the Tools Git tab only when the current project root is exactly a Git
  working-tree root and the host has a `git` executable. Accept both `.git`
  directories and worktree `.git` files. The Worktree menu lists all worktrees
  registered with the current repository, sorted by path, and marks the active
  worktree. Selecting one switches Notes' project context while preserving the
  active terminal pane and its other Notes context roots.
- Provide a worktree creation modal with a branch picker listing local branches
  alphabetically and an option to create a branch. Default the destination to
  `$HOME/Documents/ttyphoon/.git_worktrees/<sha256>`, hashing the remote origin
  URL or, when no origin is configured, the canonical project root. Accept
  relative paths within the project and absolute destinations; validate branch
  names, reject symlinked destination parents, and invoke Git with
  argument-separated commands. Add an in-repository destination to the shared
  repository `info/exclude` so a nested worktree does not appear untracked.
- Start creation from Notes asynchronously. Close the modal immediately, create
  a sticky `Creating worktree` terminal notification before invoking Git, and
  close that notification when the Git operation completes. Report failures as
  terminal error notifications.
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
- Show Git path and status-code hover help through Notes' shared tooltip rather
  than native webview tooltips. Keep tooltip behavior delegated from the Git
  pane so dynamically rendered status rows use the same path.
- Handle filenames as NUL-delimited Git status records, including spaces and
  rename source paths. Cap text diff inputs and reject path traversal and
  symlink escapes.

## Consequences

- The panel follows the active project context. A root change refreshes status;
  a non-Git root hides the tab.
- The top-level Notes Diff tab is available when the active file appears in the
  cached Git status, even when that file was opened from the Notes file list.
  Watcher-published status updates can show or hide the tab as the file becomes
  dirty or clean.
- Selecting a Git file makes it the active Notes document. The normal Notes
  filename, recent list, and history apply; the diff overlays the regular editor
  while keeping its model synchronized. Leaving the diff reveals the regular
  editor for the same file. The Diff tab is present only while a Git diff is
  active. Cmd/Ctrl+S keeps the normal format-and-save behavior.
- If Monaco is disabled, opening a Git file still loads it into Notes and
  updates history/recents, but the side-by-side diff is unavailable.
- Git operations include read, stage, unstage, worktree listing, local branch
  listing, worktree creation, and switching Notes to a registered worktree.
  Commit, branch deletion/renaming, multi-repository discovery, and persistent
  workspace lists are out of scope.

## Reference

- `utils/git/git.go` - repository validation, status, diff, stage/unstage, and
  worktree/branch operations
- `frontend.go` - Wails wrappers
- `window/backend/renderer_webkit/notes.go` - Notes project-context worktree switch
- `frontend/src/notes-git-panel.js` - status panel and file actions
- `frontend/src/shared_tooltip.js` - shared tooltip rendering and pointer tracking
- `frontend/src/monaco_adapter.js` - shared editor hooks and diff-editor facade
- `frontend/src/notes.js`, `frontend/src/notes.css` - tab, workspace, and editor
  integration
- `utils/git/git_test.go`, `frontend/src/notes-git-panel.test.js` - behavior
  coverage
- ADR 0013 - Wails binding checklist
- ADR 0018 - Monaco provider architecture
