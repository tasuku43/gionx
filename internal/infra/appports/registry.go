package appports

import (
	"time"

	"github.com/tasuku43/kra/internal/stateregistry"
)

func TouchStateRegistry(root string) error {
	return stateregistry.Touch(root, time.Now())
}
