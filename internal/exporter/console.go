package exporter

import (
	"os"

	"github.com/marianozunino/goq/internal/config"
	"github.com/streadway/amqp"
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

func (w *ConsoleExporter) WriteMessage(msg amqp.Delivery) error {
	output, err := writeMessageCommon(msg, w.config.PrettyPrint)
	if err != nil {
		return err
	}

	// Write to stdout
	_, err = os.Stdout.Write(output)
	return err
}

func (w *ConsoleExporter) Close() error {
	// No-op for console writer
	return nil
}

