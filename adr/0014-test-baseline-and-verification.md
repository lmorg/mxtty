# 14. Known test baseline and verification commands

Date: 2026-09-06

## Status

Accepted. Updated 2026-10-14: the historical failures below remain resolved,
but current verification has the unrelated failures recorded below.

## Context

`go test ./...` did not pass cleanly, and one failure *panicked*, aborting the
test binary and hiding every result that would have followed it in that package.
Without a recorded baseline each session re-diagnosed the same failures and risked
attributing them to work in progress.

## Decision

Record the baseline and the verification commands, and keep this record current as
failures are fixed.

### Current known failure

Observed on macOS with Go 1.27.1 on 2026-10-14:

- `ai/tools/file`: `TestInsertLinesCallDescriptionRendersEscapedInserts` and
  `TestPatchFileCallDescriptionRendersEmbeddedTemplate` expect older description
  formatting/line numbers.
- `ai/tools/shell`: `TestBashReturnsOutputAndExitError` expects a returned Bash
  exit error, but receives nil.
- `utils/git`: `TestRenderWorktreeInfoKeepsMultilineDiffInOneTableRow` expects
  the older two-column diff table instead of the current details-based output.

All four Go failures also reproduce in an isolated clean `HEAD` export, without
the Notes navigation/storage changes. `go build ./...` passes. The affected
Notes/history, shared path, asset, image, graph and session tests pass.

The current working-tree frontend run has 481 passes and two Git-panel failures
in `src/notes.test.js`: `resets Commits on a workspace change even when the
project root stays the same` (commit-content CSS expectations), and `resets the
branch selection when Create branch is cancelled and reports worktree action
errors` (worktree creation/removal error expectations). Git-panel/CSS edits
were already present before the Notes navigation work and are not changed by
that task. Notes rename/delete/grep, structured Edit Find and shortcut regression
tests pass.

Terminal hover-exit verification also passes: `go test ./virtualterm
./window/backend/renderer_webkit ./tmux ./types . -count=1`, `go build ./...`
and all 10 tests in `src/terminal.test.js`. The frontend-wide failures above are
unchanged.

### Resolved

Both fixed at root cause rather than by adjusting assertions:

1. **`ai/agent` panic** — `findService` called `panic()` for an unknown service.
   Every caller already handled a `nil` return, and
   `SetServiceModelFromSelection` even contained unreachable `if service == nil`
   error handling. The panic contradicted the design of its own callers, and made
   any `Agent` constructed without loaded config unusable. `findService` now
   returns `nil`; `SwitchServiceModel` gained the nil/bounds guard it lacked.

2. **`ai/tools/subagents`** — `TestSubagentToolContracts` asserted
   `strings.Contains(desc, "written summary")`, but the embedded description is
   hard-wrapped markdown containing `"written\nsummary"`. The assertion could
   never match. Comparison now collapses whitespace via `strings.Fields`, so
   reflowing a description no longer breaks the test.

### Commands

```bash
go build ./...
go test ./... -count=1                 # whole module
go test ./ai/... -count=1              # AI subsystem
go test ./ai/agent -count=1 -race      # concurrency-sensitive packages
go test ./ai/agent -run '^$'           # compile/type-check only, runs nothing

cd frontend
npm test -- --run src/notes.test.js --reporter=dot
node --check src/notes.js                            # syntax only
```

Use `npm test` rather than a bare `npx vitest`, which may prompt to install a
different major version.

## Consequences

- The AI packages now pass end to end, so a new failure there is genuinely caused
  by the change in progress.
- **Lesson: prefer fixing the code over relaxing the assertion.** Both failures
  were real defects — a panic that broke testability and contradicted its
  callers, and an assertion that could never pass — not stale expectations.
- **Lesson: a panic in a shared lookup is hostile to testing.** Returning `nil`
  and letting callers decide is both more testable and more robust at runtime.
- `-race` remains worthwhile for the emitter, agent tool maps, and `sessiondb`;
  those concurrency bugs do not reproduce reliably without it.

## Reference

- `ai/agent/services.go` — `findService`, `SwitchServiceModel`
- `ai/tools/subagents/subagent_test.go` — `containsPhrase`
- `utils/runewidth/runewidth_test.go` — passes in the 2026-10-14 verification
