package cli

import "github.com/tasuku43/kra/internal/infra/appports"

func (c *CLI) touchStateRegistry(root string) error {
	return appports.TouchStateRegistry(root)
}
