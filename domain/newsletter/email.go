package newsletter

import (
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Email consists of fields which describe an email.
type Email struct {
	ID    id.Email    `json:"id"`
	Name  string      `json:"name"`
	Email types.Email `json:"email"`
}
