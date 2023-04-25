//go:build dependency

package dependency

// This import ensures that the dependency remains in the go.mod.
import (
	_ "go.strv.io/tea/cmd/tea"
)
