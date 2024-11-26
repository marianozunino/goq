package exporter

import (
	"fmt"
	"os"

	"github.com/marianozunino/goq/internal/config"
	"github.com/streadway/amqp"
)

type FileExporter struct {
	file   *os.File
	config *config.Config
}

var _ Exporter = &FileExporter{}

func NewFileWriter(cfg *config.Config) (*FileExporter, error) {
	var file *os.File
	var err error

	// Open file based on config
	switch cfg.FileMode {
	case "append":
		file, err = os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	case "overwrite", "":
		file, err = os.Create(cfg.OutputFile)
	default:
		return nil, fmt.Errorf("invalid file mode: %s (use 'append' or 'overwrite')", cfg.FileMode)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open/create output file: %v", err)
	}
	return &FileExporter{
		file:   file,
		config: cfg,
	}, nil
}

func (w *FileExporter) WriteMessage(msg amqp.Delivery) error {
	output, err := writeMessageCommon(msg, w.config.PrettyPrint)
	if err != nil {
		return err
	}

	// Write to file
	_, err = w.file.Write(output)
	return err
}

func (w *FileExporter) Close() error {
	return w.file.Close()
}

