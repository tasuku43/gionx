package stateregistry

import (
	"time"

	base "github.com/tasuku43/gionx/internal/stateregistry"
)

type Entry = base.Entry

func Path() (string, error)                  { return base.Path() }
func Load(path string) ([]Entry, error)      { return base.Load(path) }
func Touch(root string, now time.Time) error { return base.Touch(root, now) }
