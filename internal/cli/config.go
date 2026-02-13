package cli

import (
	"fmt"

	"github.com/tasuku43/kra/internal/config"
	"github.com/tasuku43/kra/internal/infra/paths"
)

func (c *CLI) loadMergedConfig(root string) (config.Config, error) {
	globalPath, err := paths.ConfigPath()
	if err != nil {
		return config.Config{}, fmt.Errorf("resolve global config path: %w", err)
	}
	globalCfg, err := config.LoadFile(globalPath)
	if err != nil {
		return config.Config{}, fmt.Errorf("load global config %s: %w", globalPath, err)
	}

	rootPath := paths.RootConfigPath(root)
	rootCfg, err := config.LoadFile(rootPath)
	if err != nil {
		return config.Config{}, fmt.Errorf("load root config %s: %w", rootPath, err)
	}

	merged := config.Merge(globalCfg, rootCfg)
	if err := merged.Validate(); err != nil {
		return config.Config{}, err
	}
	return merged, nil
}
