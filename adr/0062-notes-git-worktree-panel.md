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
- While the Git Tools tab is selected, watch the per-worktree Git metadata
  directory and the shared Git common directory resolved by `git rev-parse`.
  Debounce metadata changes before refreshing status; stop both watches when the
  Git tab is deselected. Watch metadata writes as well as entry creation/removal,
  but exclude object, hook, and reflog trees. Ordinary working-tree file changes
  continue to refresh status through the Notes file-index watcher.
- Run panel status reads with `git --no-optional-locks status` so Git does not
  persist optional index refreshes and retrigger its own metadata watcher.
  Stage/unstage and worktree mutations retain their normal locking behavior.
- Enable the Tools Git tab only when the current project root is exactly a Git
  working-tree root and the host has a `git` executable. Accept both `.git`
  directories and worktree `.git` files. The Worktree menu lists all worktrees
  registered with the current repository, sorted by path, and marks the active
  worktree. Selecting one switches Notes' project context while preserving the
  active terminal pane and its other Notes context roots.
- Place a `Branches` button after `Commits`. Its menu starts with Create new
  branch and a separator, followed by sorted local branches. Creation uses the
  shared text prompt and creates/checks out the branch. A branch submenu offers
  copy name, checkout, Fetch latest, and Create new worktree. Checkout reports
  `Switched to <branch>` on success and Git errors on failure without forcing.
  Fetch targets the configured upstream (or same-named branch on origin), never
  checks out or pulls, and reports stdout as info and stderr as error. Worktree
  creation uses the default directory plus branch name and the shared creation/
  switch flow. Git commands remain in Go behind thin Wails bindings.
- Provide a worktree creation modal with a branch picker listing local branches
  alphabetically and an option to create a branch. Default the destination to
  `$HOME/Documents/ttyphoon/.git_worktrees/<sha256>/<branch>`, hashing the
  remote origin URL or, when no origin is configured, the canonical project
  root. Accept relative paths within the project and absolute destinations;
  validate branch names, reject symlinked destination parents, and invoke Git
  with argument-separated commands. Add an in-repository destination to the
  shared repository `info/exclude` so a nested worktree does not appear
  untracked.
- Start creation from Notes asynchronously. Close the modal immediately, create
  a sticky `Creating worktree` terminal notification before invoking Git, and
  invoke the Worktrees menu's Switch action after successful creation, using the
  returned canonical path. Both dialog and commit-menu creation use this shared
  flow; close the notification after switching completes. Report failures as
  terminal error notifications.
- Removing a worktree removes only its working directory/registration and never
  deletes its branch. Do not force removal; let Git reject dirty worktrees.
- Render branch state, ahead/behind counts, and changed paths in the Git pane.
  The `Commits` toggle is the first right-aligned Git toolbar item, before
  Worktrees. It replaces the left sidebar file tree
  with topologically ordered history, while Tools continues to show Git status.
  Keep the shared file/grep filter visible during history mode. Use the same
  `FilterStrings` query language against full commit hashes and complete messages,
  including message bodies, while still displaying short hashes and subjects.
  A non-empty history filter hides the graph and uses the full sidebar width for
  commit text. Clearing it restores the unfiltered graph and text positioning.
  Search a cached window of up to 10,000 commits; repeated queries reuse it,
  with request guards rejecting stale filter results. Restore file/grep filtering
  on history exit; keep the file tree mounted and updated in the background.
  History includes only commits reachable from the checked-out worktree's HEAD,
  including merged ancestry, not unrelated branches, remote refs, or tags.
  Detached HEAD uses its own ancestry; an unborn branch shows no commits.
  Go reads Git's graph layout and
  NUL-delimited short hash, first-line subject, author date/time, and author name;
  display the parenthesized author on a separate line below the date/time.
  the frontend draws the graph with the project tree's one-pixel lines in theme
  yellow. Only the portion behind the text is masked to 0.7 alpha; hovering text
  dims its graph to 0.2 opacity, while hovering exposed graph dims its text by
  the same factor. A row resize observer keeps the overlap mask aligned and is
  disconnected when history is replaced, hidden, or its project changes. The renderer
  connects lanes using adjacent graph lines rather than fixed slash offsets.
  Adjacent commit SVGs share boundary coordinates, so expanding/collapsing lanes
  remain continuous even when commit text makes rows different heights. Merge
  dots connect to subsequent lanes instead of being drawn as horizontal dashes.
  Opposite-direction slashes in the same column retain that lane. Leftward
  collapse connectors are kept even where a vertical continuation also exists;
  tests verify lane coverage, not just matching endpoints between commit rows.
  Same-direction diagonal continuations take precedence over adjacent vertical
  rails. Route diagonals first so crossing rails cannot create a second, false
  input to a slash that already continues from another slash.
  The graph is anchored to the left while commit text is anchored to the right. Wide graphs
  may overlap text; text is painted above the graph and retains up to 70% of the
  list width (or its natural width, if smaller), using more space when available.
  Commit rows use the grep-result hover style. Commit subjects use the Markdown paragraph
  font (Lato) at the Notes minus-two
  font size. Load 100 commits initially, with explicit
  Load more batches up to 10,000. Turning the toggle off restores files. Workspace
  changes reset the toggle, including workspaces sharing the same project root.
  A left-side comparison selector defaults to `HEAD`; local branches are listed
  alphabetically with the checked-out branch omitted. Selecting another branch
  changes only the diff base, never the checked-out branch. In that mode the
  file list contains paths differing from the selected branch plus untracked
  working-tree files, and each file diff uses that branch as its original side.
  Clicking a history commit (or pressing Enter/Space on its row) selects that
  commit through the same comparison handler as the branch menu. The selector
  shows its short hash; comparison APIs use its validated full commit ID.
  This refreshes the file list and diff base without checking out the commit.
  Right-clicking a history row opens a menu titled with its short hash: Copy
  commit hash, separator, Compare commit, Create worktree. Copy uses the full
  hash; Compare uses icon `0xe13a` and invokes the left-click handler. Create worktree targets the
  commit in detached mode at the default worktree directory plus its short hash,
  while Git still receives the full hash. Detached worktree menu labels use
  Git's short abbreviation of their HEAD instead of `detached HEAD`,
  reusing the dialog's shared notification-and-creation handler and Wails binding.
  Show a provider-brand button before Refresh only for GitHub, Azure DevOps,
  GitLab, or Bitbucket origins. Resolve HTTPS browser URLs in Go using the
  origin and current branch; open them through the system browser. Hide the
  button when origin is missing, unsupported, or the worktree is detached.
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
- Branch comparison is view-only: `HEAD` remains the default status view, while
  selecting a local branch changes the file list and diff base without checking
  out or mutating any branch.
- Git metadata status refreshes are watched only while the Git tab is selected.
  Linked worktrees watch their external per-worktree and common Git directories;
  ordinary file edits are covered by the active project file-index watch.
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
