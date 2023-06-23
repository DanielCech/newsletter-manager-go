package query

import _ "embed"

var (
	//go:embed scripts/Create.sql
	Create string
	//go:embed scripts/Update.sql
	Update string
	//go:embed scripts/Read.sql
	Read string
	//go:embed scripts/ListByAuthor.sql
	ListByAuthor string
	//go:embed scripts/List.sql
	List string
)
