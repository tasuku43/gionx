package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *CLI) runShell(args []string) int {
	if len(args) == 0 {
		c.printShellUsage(c.Err)
		return exitUsage
	}

	switch args[0] {
	case "-h", "--help", "help":
		c.printShellUsage(c.Out)
		return exitOK
	case "init":
		return c.runShellInit(args[1:])
	default:
		fmt.Fprintf(c.Err, "unknown command: %q\n", strings.Join(append([]string{"shell"}, args[0]), " "))
		c.printShellUsage(c.Err)
		return exitUsage
	}
}

func (c *CLI) runShellInit(args []string) int {
	if len(args) > 1 {
		fmt.Fprintf(c.Err, "unexpected args for shell init: %q\n", strings.Join(args[1:], " "))
		c.printShellUsage(c.Err)
		return exitUsage
	}

	shellName := ""
	if len(args) == 1 {
		shellName = strings.TrimSpace(args[0])
	} else {
		shellName = detectShellName()
	}
	if shellName == "" {
		shellName = "zsh"
	}

	script, err := renderShellInitScript(shellName)
	if err != nil {
		fmt.Fprintf(c.Err, "render shell init script: %v\n", err)
		return exitUsage
	}
	fmt.Fprint(c.Out, script)
	return exitOK
}

func detectShellName() string {
	raw := strings.TrimSpace(os.Getenv("SHELL"))
	if raw == "" {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(filepath.Base(raw)))
}

func renderShellInitScript(shellName string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(shellName)) {
	case "zsh", "bash", "sh":
		return renderPOSIXShellInitScript(shellName), nil
	case "fish":
		return renderFishShellInitScript(), nil
	default:
		return "", fmt.Errorf("unsupported shell %q (supported: zsh, bash, sh, fish)", shellName)
	}
}

func renderPOSIXShellInitScript(shellName string) string {
	return fmt.Sprintf(`# kra shell integration (%s)
# Add this line to your shell rc file (~/.%src), then restart shell:
#   eval "$(kra shell init %s)"
kra() {
  local __kra_action_file __kra_status
  __kra_action_file="$(mktemp "${TMPDIR:-/tmp}/kra-action.XXXXXX")" || return 1
  KRA_SHELL_ACTION_FILE="$__kra_action_file" command kra "$@"
  __kra_status=$?
  if [ $__kra_status -ne 0 ]; then
    rm -f "$__kra_action_file"
    return $__kra_status
  fi
  if [ -s "$__kra_action_file" ]; then
    eval "$(cat "$__kra_action_file")"
  fi
  rm -f "$__kra_action_file"
}
`, shellName, shellName, shellName)
}

func renderFishShellInitScript() string {
	return `# kra shell integration (fish)
# Add this line to your config.fish, then restart shell:
#   eval (kra shell init fish)
function kra
  set -l __kra_action_file (mktemp "/tmp/kra-action.XXXXXX"); or return 1
  env KRA_SHELL_ACTION_FILE="$__kra_action_file" command kra $argv
  set -l __kra_status $status
  if test $__kra_status -ne 0
    rm -f "$__kra_action_file"
    return $__kra_status
  end
  if test -s "$__kra_action_file"
    eval (cat "$__kra_action_file")
  end
  rm -f "$__kra_action_file"
end
`
}
