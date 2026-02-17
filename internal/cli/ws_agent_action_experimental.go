//go:build experimental

package cli

func (c *CLI) wsAgentActionEnabled() bool {
	return true
}

func (c *CLI) runWSRunAgent(args []string) int {
	return c.runAgentRun(args)
}
