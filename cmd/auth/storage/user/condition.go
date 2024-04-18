package user

type Condition interface {
	condition()
}

// SearchEmail is a condition to search for a user by email.
type SearchEmail struct {
	Email string
}

func (SearchEmail) condition() {}
