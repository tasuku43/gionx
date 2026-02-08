package paths

import base "github.com/tasuku43/gionx/internal/paths"

func ResolveExistingRoot(cwd string) (string, error) { return base.ResolveExistingRoot(cwd) }
func StateDBPathForRoot(root string) (string, error) { return base.StateDBPathForRoot(root) }
func DefaultRepoPoolPath() (string, error)           { return base.DefaultRepoPoolPath() }
func WriteCurrentContext(root string) error          { return base.WriteCurrentContext(root) }
