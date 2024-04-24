package organization

type Condition interface {
	condition()
}

// SearchEmail is a condition to search for a user by email.
type SearchDisplayName struct {
	DisplayName string
}

func (SearchDisplayName) condition() {}
