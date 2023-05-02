package newsletter

// Subscription consists of fields which describe an subscription.
type Subscription struct {
	Email        string `json:"email"`
	NewsletterID int    `json:"newsletterId"`
}

type FullSubscription struct {
	Email        string `json:"email"`
	NewsletterID int    `json:"newsletterId"`
	Token        string `json:"token"`
}
