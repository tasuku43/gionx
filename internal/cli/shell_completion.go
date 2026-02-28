package cli

import (
	"fmt"
	"strings"
)

var kraCompletionRootCommands = []string{
	"init",
	"bootstrap",
	"context",
	"repo",
	"template",
	"shell",
	"ws",
	"doctor",
	"version",
	"help",
}

var kraCompletionGlobalFlags = []string{
	"--debug",
	"--version",
	"--help",
	"-h",
}

var kraCompletionSubcommandOrder = []string{
	"bootstrap",
	"context",
	"repo",
	"template",
	"shell",
	"ws",
}

var kraCompletionSubcommands = map[string][]string{
	"bootstrap": {"agent-skills", "help"},
	"context":   {"current", "list", "create", "use", "rename", "rm", "help"},
	"repo":      {"add", "discover", "remove", "gc", "help"},
	"template":  {"validate", "help"},
	"shell":     {"init", "completion", "help"},
	"ws": {
		"create",
		"import",
		"list",
		"ls",
		"dashboard",
		"insight",
		"select",
		"lock",
		"unlock",
		"open",
		"switch",
		"add-repo",
		"remove-repo",
		"close",
		"reopen",
		"purge",
		"help",
	},
}

func renderShellCompletionScript(shellName string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(shellName)) {
	case "zsh":
		return renderZshCompletionScript(), nil
	case "bash", "sh":
		return renderBashCompletionScript(), nil
	case "fish":
		return renderFishCompletionScript(), nil
	default:
		return "", fmt.Errorf("unsupported shell %q (supported: zsh, bash, sh, fish)", shellName)
	}
}

func renderBashCompletionScript() string {
	return fmt.Sprintf(`# kra completion (bash)
_kra_completion() {
  local cur prev cmd i
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev=""
  if [[ ${COMP_CWORD} -gt 0 ]]; then
    prev="${COMP_WORDS[COMP_CWORD-1]}"
  fi

  cmd=""
  for ((i=1; i<COMP_CWORD; i++)); do
    if [[ "${COMP_WORDS[i]}" != -* ]]; then
      cmd="${COMP_WORDS[i]}"
      break
    fi
  done

  if [[ -z "${cmd}" ]]; then
    COMPREPLY=( $(compgen -W "%s" -- "${cur}") )
    return 0
  fi

  case "${cmd}" in
%s
  esac

  return 0
}
complete -o default -F _kra_completion kra
`, strings.Join(kraCompletionTopWords(), " "), renderBashSubcommandCases())
}

func renderBashSubcommandCases() string {
	lines := make([]string, 0, len(kraCompletionSubcommandOrder)*4)
	for _, cmd := range kraCompletionSubcommandOrder {
		subs := strings.Join(kraCompletionSubcommands[cmd], " ")
		lines = append(lines,
			fmt.Sprintf("    %s)", cmd),
			fmt.Sprintf("      if [[ \"${prev}\" == \"%s\" ]]; then", cmd),
			fmt.Sprintf("        COMPREPLY=( $(compgen -W \"%s\" -- \"${cur}\") )", subs),
			"      fi",
			"      ;;",
		)
	}
	return strings.Join(lines, "\n")
}

func renderZshCompletionScript() string {
	return fmt.Sprintf(`# kra completion (zsh)
_kra_completion() {
  local -a top sub
  local cmd="" i

  top=(%s)

  for (( i=2; i<CURRENT; i++ )); do
    if [[ "${words[i]}" != -* ]]; then
      cmd="${words[i]}"
      break
    fi
  done

  if [[ -z "${cmd}" ]]; then
    compadd -- $top
    return 0
  fi

  sub=()
  case "$cmd" in
%s
  esac

  if [[ ${#sub[@]} -gt 0 ]] && [[ "${words[CURRENT-1]}" == "$cmd" ]]; then
    compadd -- $sub
  fi
}
compdef _kra_completion kra
`, strings.Join(kraCompletionTopWords(), " "), renderZshSubcommandCases())
}

func renderZshSubcommandCases() string {
	lines := make([]string, 0, len(kraCompletionSubcommandOrder))
	for _, cmd := range kraCompletionSubcommandOrder {
		subs := strings.Join(kraCompletionSubcommands[cmd], " ")
		lines = append(lines, fmt.Sprintf("    %s) sub=(%s) ;;", cmd, subs))
	}
	return strings.Join(lines, "\n")
}

func renderFishCompletionScript() string {
	var b strings.Builder
	b.WriteString("# kra completion (fish)\n")
	b.WriteString("complete -c kra -f\n")
	b.WriteString("complete -c kra -l debug -d \"Enable debug logging\"\n")
	b.WriteString("complete -c kra -l version -d \"Print version and exit\"\n")
	b.WriteString("complete -c kra -l help -s h -d \"Show help\"\n")
	b.WriteString(
		fmt.Sprintf(
			"complete -c kra -n \"__fish_use_subcommand\" -a \"%s\"\n",
			strings.Join(kraCompletionRootCommands, " "),
		),
	)
	for _, cmd := range kraCompletionSubcommandOrder {
		b.WriteString(
			fmt.Sprintf(
				"complete -c kra -n \"__fish_seen_subcommand_from %s\" -a \"%s\"\n",
				cmd,
				strings.Join(kraCompletionSubcommands[cmd], " "),
			),
		)
	}
	return b.String()
}

func kraCompletionTopWords() []string {
	out := make([]string, 0, len(kraCompletionRootCommands)+len(kraCompletionGlobalFlags))
	out = append(out, kraCompletionRootCommands...)
	out = append(out, kraCompletionGlobalFlags...)
	return out
}
