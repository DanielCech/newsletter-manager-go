package query

import _ "embed"

var (
	//go:embed scripts/CreateRefreshToken.sql
	CreateRefreshToken string
	//go:embed scripts/ReadRefreshToken.sql
	ReadRefreshToken string
	//go:embed scripts/DeleteRefreshToken.sql
	DeleteRefreshToken string
	//go:embed scripts/DeleteRefreshTokensByAuthorID.sql
	DeleteRefreshTokensByAuthorID string
)
