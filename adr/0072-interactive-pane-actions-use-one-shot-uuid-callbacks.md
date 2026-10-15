# 72. Interactive pane actions use one-shot UUID callbacks

Date: 2026-10-15

## Status

Accepted.

## Context

Some commands must run in an interactive cloned tmux pane so shell hooks,
credential prompts, and other terminal behavior work normally. Their caller
still needs the command's screen output and exit status. A completion escape
sequence must not be actionable merely because it appears in copied or replayed
terminal text; it must correspond to an outstanding invocation.

## Decision

- Before starting an interactive pane action, register its callback and generate
  a random UUIDv4 for that invocation. Pass that UUID into the generated shell
  command so its completion hook can include it in a ttyphoon APC. Start the
  wrapper as the new pane's command with tmux `split-window -k`; the pane remains
  available after the command exits. The wrapper preserves the wrapped command's
  exit status even when the shell uses `errexit`.
- Store callbacks in a mutex-protected registry owned by the relevant `Tmux`
  instance, keyed by UUID. Do not use a package-global callback map.
- Accept a completion APC only when its UUID matches a pending entry. Remove the
  entry atomically before invoking its callback, so each UUID can complete at
  most once. Unknown, stale, malformed, or duplicate completion APCs are ignored.
- The command may finish and emit its APC before tmux has initialized the new
  pane in ttyphoon. Queue that result by UUID and source pane ID; bind and deliver
  it only after pane initialization and renderer refresh complete.
- The callback receives the terminal screen output captured through the
  existing screen-reading API and the script exit code. Callback execution
  occurs outside the registry lock.
- Bash and Zsh output-block integrations report their exported environment as
  base64-encoded NUL-delimited data, allowing a cloned pane to inherit the
  focused shell's reported environment as well as its working directory.
- Remove a pending callback if pane creation fails or the associated pane exits
  without reporting completion.
- On exit code zero, paint the pane with the existing green hover chrome and
  close it after thirty seconds. On non-zero exit, paint it with the existing red
  error-hover chrome and leave it open for inspection.
- Include the backend auto-close deadline as Unix milliseconds in the pane
  snapshot's `autoCloseAt` field. Show the notification-style pie countdown next
  to the close icon, using that deadline so redraws and tab switches do not reset
  it. Animate the DOM indicator only, without terminal redraws or tmux commands.
  Clicking the timer calls `TerminalCancelPaneAutoClose`, atomically cancelling
  the deadline in Go and removing the indicator. The pane then stays open until
  manually closed with its close icon. Failed panes have no countdown. Respawned
  panes clear old deadlines, and expired callbacks verify pane identity and
  atomically claim their original deadline before closing.
- Publish changed visible-pane snapshots to the frontend as `terminalPanes`
  events, with pane ID, `x`, `y`, `width`, and `height` in terminal cells and a
  `canClose` flag. `TerminalGetPanes` supplies the initial snapshot. Refresh
  geometry on tmux `%layout-change` notifications as well as terminal resizes,
  tab switches, and zoom changes; the redraw loop also detects completion-state
  changes without polling tmux per frame.
- Geometry synchronization must not echo tmux-reported dimensions back as a
  `resize-pane` command. `PaneT.Resize` returns without a command when the size
  already matches the known tmux geometry. Otherwise layout notifications form
  a resize/refresh feedback loop, causing terminal flicker and application
  lockups. Coalesce layout-event bursts into one refresh worker with one pending
  trailing refresh rather than starting an unbounded goroutine per event.
- Place a close icon inside the top-right corner of completed panes. Both
  validated interactive-action results and tmux's `pane_dead` state make a pane
  closable; a running shell is not closable through this control. Track the
  `pane-died` hook for retained panes and clear completion state on respawn.
  `TerminalClosePane` delegates to tmux, which validates the stored pane state,
  kills that exact pane ID, and reconciles removed panes before refreshing the
  renderer. Manual close remains available during and after the countdown;
  closing the pane makes its pending auto-close callback a no-op.
- Keep this terminal mechanism generic; Git-specific execution policy remains
  in `utils/git` (ADR 0071).
- Markdown code blocks and Jupyter's Send to terminal action share the
  `SendToTerminal(content)` Wails API. It delegates to `CloneActivePane`, running
  the supplied content immediately as a shell script in a cloned tmux pane,
  rather than pasting into the active shell. The clone inherits the source
  pane's environment and working directory, uses the same authenticated
  completion callback, and gets the same close controls and cancellable success
  countdown. Pane setup errors are returned to the frontend; tmux must be
  available rather than falling back to writing into the original pane.

The UUID match is an accidental-use guard, not a cryptographic security
boundary.

## Consequences

Interactive callers receive a one-shot completion result without polling pane
text or running the command a second time. Callback lifetime follows the tmux
session, and cleanup prevents callbacks from accumulating for panes that never
complete normally.

## Reference

- `tmux/pane.go` - cloned pane lifecycle
- `tmux/bash.sh.tmpl` - completion APC wrapper
- `integrations/shell.bash`, `integrations/shell.zsh` - command environment snapshots
- `virtualterm/ansi_apc.go` - APC parsing entry point
- `virtualterm/term.go` - terminal screen output access
- `utils/git/action.go` - Git execution policy and interactive action caller
- `frontend.go` - Wails-to-renderer adapter