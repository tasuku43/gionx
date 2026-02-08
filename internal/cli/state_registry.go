package cli

import "github.com/tasuku43/gionx/internal/infra/appports"

func (c *CLI) touchStateRegistry(root string) error {
	return appports.TouchStateRegistry(root)
}
