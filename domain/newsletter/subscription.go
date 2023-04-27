package newsletter

import (
	"newsletter-manager-go/types"
	"newsletter-manager-go/types/id"
)

// Subscription consists of fields which describe an subscription.
type Subscription struct {
	ID    id.Subscription `json:"id"`
	Name  string          `json:"name"`
	Email types.Email     `json:"email"`
}
