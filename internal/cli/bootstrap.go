package cli

import (
	"fmt"
	"strings"
)

func (c *CLI) runBootstrap(args []string) int {
	if len(args) == 0 {
		c.printBootstrapUsage(c.Err)
		return exitUsage
	}

	switch args[0] {
	case "-h", "--help", "help":
		c.printBootstrapUsage(c.Out)
		return exitOK
	case "agent-skills":
		return c.runBootstrapAgentSkills(args[1:])
	default:
		fmt.Fprintf(c.Err, "unknown command: %q\n", strings.Join(append([]string{"bootstrap"}, args[0]), " "))
		c.printBootstrapUsage(c.Err)
		return exitUsage
	}
}
