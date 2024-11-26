package app

import "github.com/marianozunino/goq/internal/config"

// Monitor is a package-level function for convenience
func Monitor(cfg *config.Config) error {
	processor, err := NewMessageProcessor(cfg)
	if err != nil {
		return err
	}
	return processor.Monitor()
}
