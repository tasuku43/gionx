//go:build !experimental

package cli

func (c *CLI) wsAgentActionEnabled() bool {
	return false
}

func (c *CLI) runWSRunAgent(_ []string) int {
	return exitUsage
}
