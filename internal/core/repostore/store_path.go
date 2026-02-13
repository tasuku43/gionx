package repostore

import (
	"path/filepath"

	"github.com/tasuku43/gionx/internal/core/repospec"
)

func StorePath(bareRoot string, spec repospec.Spec) string {
	return filepath.Join(bareRoot, spec.Host, spec.Owner, spec.Repo+".git")
}
