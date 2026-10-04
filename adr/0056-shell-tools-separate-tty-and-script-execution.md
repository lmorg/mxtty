# 0056 - Separate interactive terminal commands from Bash scripts

Date: 2026-10-04

## Status

Accepted.

## Context

A command sent to the active shell's TTY and a script executed by Bash have
fundamentally different execution semantics. Interactive programs need the live
terminal, while pipelines, redirects, shell variables, and control flow should
run as a script without a TTY. Routing both through the terminal made it
unclear which shell syntax was valid and prevented reliable output streaming
for scripted work.

## Decision

- Replace the `commandLine` tool with `interactive`. It sends one command
  invocation directly to the active terminal for programs that need a TTY. Its
  input must not contain Bash or POSIX shell syntax; callers use `bash` for
  scripts and composed shell commands.
- Add a `bash` tool that accepts a multiline script, writes it to a uniquely
  named temporary directory, and runs `bash <script>` without a TTY.
- Stream combined stdout and stderr through one typed `tool-output` block and
  return non-zero exit errors with captured output to the model for correction.
  A command failure alone does not abort the overall agent run.
- Remove the temporary directory after execution. Both tools use structured
  `{command}` / `{script}` inputs and require invocation permission by default.

## Consequences

- Interactive terminal state and non-interactive shell scripts are distinct API
  choices with descriptions that tell the model which syntax and behavior apply.
- Bash output appears live in the AI panel without duplicating the final output
  block.
- The shell tool does not provide a general-purpose process sandbox; the
  existing tool-permission policy remains the invocation boundary.

## References

- `ai/tools/shell/interactive.go` - terminal-backed command tool
- `ai/tools/shell/bash.go` - temporary-script execution and streamed output
- `ai/tools/shell/interactive_description.md` and `bash_description.md` - model guidance
- `ai/agent/runtime_eino.go` - tool errors returned to the model
- ADR 0009 - reported shell environment state
- ADR 0051 - structured native tool inputs
