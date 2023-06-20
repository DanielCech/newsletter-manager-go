package query

import _ "embed"

var (
	//go:embed scripts/Create.sql
	Create string
	//go:embed scripts/Update.sql
	Update string
	//go:embed scripts/Read.sql
	Read string
)
