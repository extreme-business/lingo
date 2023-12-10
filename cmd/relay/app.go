package relay

import (
	"fmt"

	"github.com/dwethmar/lingo/database"
)

// Options options for the relay server
type Options struct {
	Transactor *database.Transactor
}

// Start starts the relay server
func Start(option Options) error {
	if option.Transactor == nil {
		return fmt.Errorf("transactor is not set in options")
	}

	return nil
}
