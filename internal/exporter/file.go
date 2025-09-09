package exporter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/marianozunino/goq/internal/config"
	"github.com/wagslane/go-rabbitmq"
)

type FileExporter struct {
	writer *bufio.Writer
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
		return nil, &ExporterError{
			Type: ErrorTypeConfiguration,
			Err:  fmt.Errorf("invalid file mode: %s (use 'append' or 'overwrite')", cfg.FileMode),
		}
	}
	if err != nil {
		return nil, &ExporterError{
			Type: ErrorTypeFileIO,
			Err:  fmt.Errorf("failed to open/create output file: %v", err),
		}
	}

	writer := bufio.NewWriter(file)

	return &FileExporter{
		writer: writer,
		file:   file,
		config: cfg,
	}, nil
}

func (w *FileExporter) WriteMessage(msg rabbitmq.Delivery) error {
	output, err := writeMessageCommon(msg, w.config.PrettyPrint)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(output)
	if err != nil {
		return &ExporterError{
			Type: ErrorTypeFileIO,
			Err:  fmt.Errorf("failed to write to file: %v", err),
		}
	}

	if err := w.writer.Flush(); err != nil {
		return &ExporterError{
			Type: ErrorTypeFileIO,
			Err:  fmt.Errorf("failed to flush buffer: %v", err),
		}
	}

	return nil
}

func (w *FileExporter) Close() error {
	if err := w.writer.Flush(); err != nil {
		return &ExporterError{
			Type: ErrorTypeFileIO,
			Err:  fmt.Errorf("failed to flush buffer on close: %v", err),
		}
	}
	if err := w.file.Close(); err != nil {
		return &ExporterError{
			Type: ErrorTypeFileIO,
			Err:  fmt.Errorf("failed to close file: %v", err),
		}
	}
	return nil
}
