Interactive terminal tool. Sends the provided command directly to the active terminal as if typed by the user, for commands that require a real TTY, such as interactive programs, full-screen terminal applications, or commands that prompt for input.
- Pass a JSON object with command text in `command`, for example `{"command":"top"}`. Do not wrap the object in an `input` property.
- Provide one executable command invocation only. Do not use Bashisms or POSIX shell syntax: no pipes, redirects, command substitution, variable expansion, chaining, shell functions, or shell control structures.
- Use the `bash` tool when you need to run a script, compose commands, use pipelines, or use shell syntax.
- ANSI escape sequences are not included in the returned output.
- The returned output is the final visual state of the terminal after the command has executed.
- Commands will run on the active shell. You should take care to write command lines that conform to the programming language syntax of that shell. The shell name is included below.

## Example

Input:
```
echo "hello world" | grep "world"
```

Returns:
```json
{
	"Output":    "hello world",
	"ExitCode":  0,
	"TimeStart": "",
	"TimeEnd":   "",
	"Duration":  "1s"
}
```

