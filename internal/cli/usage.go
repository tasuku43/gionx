package cli

import (
	"fmt"
	"io"
)

func (c *CLI) printRootUsage(w io.Writer) {
	fmt.Fprint(w, `Usage:
  gionx <command> [args]

Commands:
  init              Initialize GIONX_ROOT (not implemented yet)
  ws                Workspace commands
  version           Print version
  help              Show this help

Run:
  gionx <command> --help
`)
}

func (c *CLI) printWSUsage(w io.Writer) {
	fmt.Fprint(w, `Usage:
  gionx ws <subcommand> [args]

Subcommands:
  create            Create a workspace (not implemented yet)
  list              List workspaces (not implemented yet)
  add-repo          Add repo to workspace (not implemented yet)
  close             Close workspace (not implemented yet)
  reopen            Reopen workspace (not implemented yet)
  purge             Purge workspace (not implemented yet)
  help              Show this help

Run:
  gionx ws <subcommand> --help
`)
}
