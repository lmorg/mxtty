package git

import (
	"errors"
	"fmt"
	"strings"
)

type ActionMode string

const (
	ActionModeGo          ActionMode = "go"
	ActionModeInteractive ActionMode = "interactive"
)

type InteractiveActionResult struct {
	Output   string
	ExitCode int
	Err      error
}

type InteractiveActionRunner func(command string, complete func(InteractiveActionResult)) error

type ActionRequest struct {
	Name       string
	Parameters map[string]string
}

func RunAction(root string, action ActionRequest, agentEnvironment map[string]string, interactiveRunner InteractiveActionRunner) map[string]string {
	result := map[string]string{"mode": "", "stdout": "", "stderr": "", "error": ""}
	args, err := actionArguments(action)
	if err != nil {
		result["error"] = err.Error()
		return result
	}
	mode, ok := actionMode(action.Name)
	if !ok {
		result["error"] = fmt.Sprintf("unsupported Git action %q", action.Name)
		return result
	}
	result["mode"] = string(mode)
	for _, argument := range args {
		if strings.ContainsAny(argument, "\x00\r\n") {
			result["error"] = "Git action argument contains unsupported control characters"
			return result
		}
	}

	gitPath, repoRoot, err := checkedRepository(root)
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	if mode == ActionModeInteractive {
		if interactiveRunner == nil {
			result["error"] = errors.New("interactive Git action runner is unavailable").Error()
			return result
		}
		completed := make(chan InteractiveActionResult, 1)
		if err := interactiveRunner(interactiveCommand(repoRoot, args), func(actionResult InteractiveActionResult) {
			select {
			case completed <- actionResult:
			default:
			}
		}); err != nil {
			result["error"] = err.Error()
			return result
		}
		actionResult := <-completed
		result["stdout"] = actionResult.Output
		result["exitCode"] = fmt.Sprintf("%d", actionResult.ExitCode)
		if actionResult.Err != nil {
			result["error"] = actionResult.Err.Error()
		} else if actionResult.ExitCode != 0 {
			result["error"] = fmt.Sprintf("git action failed with exit code %d", actionResult.ExitCode)
		}
		return result
	}

	output, err := runGitWithAgentEnvironment(gitPath, repoRoot, agentEnvironment, args...)
	if err != nil {
		result["error"] = err.Error()
		return result
	}
	for key, value := range output {
		result[key] = value
	}
	return result
}

func interactiveCommand(projectRoot string, args []string) string {
	parts := make([]string, 0, len(args)+3)
	parts = append(parts, "git", "-C", quoteShellArgument(projectRoot))
	for _, argument := range args {
		parts = append(parts, quoteShellArgument(argument))
	}
	return strings.Join(parts, " ")
}

func quoteShellArgument(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func actionArguments(action ActionRequest) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(action.Name)) {
	case "commit":
		message := strings.TrimSpace(action.Parameters["message"])
		if message == "" {
			return nil, errors.New("commit message is required")
		}
		return []string{"commit", "-m", message}, nil
	case "fetch":
		return []string{"fetch"}, nil
	default:
		return nil, fmt.Errorf("unsupported Git action %q", action.Name)
	}
}

func actionMode(action string) (ActionMode, bool) {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "commit":
		return ActionModeInteractive, true
	case "fetch":
		return ActionModeGo, true
	default:
		return "", false
	}
}
