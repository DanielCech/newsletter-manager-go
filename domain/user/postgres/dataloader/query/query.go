package query

import _ "embed"

// User
var (
	//go:embed scripts/ListUsersByIDs.sql
	ListUsersByIDs string
)
