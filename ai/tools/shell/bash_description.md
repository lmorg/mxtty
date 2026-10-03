Bash script tool. Writes the provided multiline script to a temporary file and executes it with Bash. Use this for non-interactive shell work, including pipelines, redirects, variables, loops, conditionals, command substitution, and other Bash or POSIX shell syntax.
- Pass a JSON object with the script in `script`, for example `{"script":"set -e\nprintf '%s\\n' hello"}`. Do not wrap the object in an `input` property.
- The script runs as a Bash script without a TTY. Use the `interactive` tool for commands that require a real terminal or interactive input.
- Standard output and standard error are returned together and streamed to the AI panel while the script runs.
