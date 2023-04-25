//go:build graphql_dependencies

package graphql

// This import ensures that the dependency remains in the go.mod.
import (
	_ "github.com/99designs/gqlgen"
)
