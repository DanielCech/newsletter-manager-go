package newsletter

import (
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Email consists of fields which describe an email.
type Email struct {
	ID    id.Email
	Name  string
	Email types.Email
}
