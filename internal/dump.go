package app

import "github.com/marianozunino/goq/internal/config"

// Dump is a package-level function for convenience
func Dump(cfg *config.Config) error {
	processor, err := NewMessageProcessor(cfg)
	if err != nil {
		return err
	}
	return processor.Dump()
}
