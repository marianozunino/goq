package filewriter

import (
	"fmt"
	"os"

	"github.com/marianozunino/goq/internal/config"
)

type Writer struct {
	file   *os.File
	config *config.Config
}

func NewWriter(cfg *config.Config) (*Writer, error) {
	var file *os.File
	var err error

	if cfg.FileMode == "append" {
		file, err = os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	} else {
		file, err = os.Create(cfg.OutputFile) // Overwrite mode
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open/create output file: %v", err)
	}

	return &Writer{
		file:   file,
		config: cfg,
	}, nil
}

func (w *Writer) WriteMessage(message string) error {
	_, err := w.file.WriteString(fmt.Sprintf("Message: %s\n", message))
	if err != nil {
		return fmt.Errorf("failed to write message to file: %v", err)
	}
	return nil
}

func (w *Writer) Close() error {
	return w.file.Close()
}
