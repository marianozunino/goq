package exporter

import (
	"fmt"
	"os"

	"github.com/marianozunino/goq/internal/config"
	"github.com/wagslane/go-rabbitmq"
)

type ConsoleExporter struct {
	config *config.Config
}

var _ Exporter = &ConsoleExporter{}

func NewConsoleExporter(cfg *config.Config) (*ConsoleExporter, error) {
	return &ConsoleExporter{
		config: cfg,
	}, nil
}

func (w *ConsoleExporter) WriteMessage(msg rabbitmq.Delivery) error {
	output, err := writeMessageCommon(msg, w.config.PrettyPrint)
	if err != nil {
		return err
	}

	// Write to stdout
	_, err = os.Stdout.Write(output)
	if err != nil {
		return &ExporterError{
			Type: ErrorTypeConsoleIO,
			Err:  fmt.Errorf("failed to write to console: %v", err),
		}
	}
	return nil
}

func (w *ConsoleExporter) Close() error {
	// No-op for console writer
	return nil
}
