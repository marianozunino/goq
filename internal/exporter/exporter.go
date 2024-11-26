package exporter

import (
	"encoding/json"
	"fmt"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/model"
	"github.com/streadway/amqp"
)

type Exporter interface {
	WriteMessage(msg amqp.Delivery) error
	Close() error
}

func NewExporter(cfg *config.Config) (Exporter, error) {
	switch cfg.Writer {
	case config.ConsoleExporterKind:
		return NewConsoleExporter(cfg)
	case config.FileWriterKind:
		return NewFileWriter(cfg)
	default:
		return nil, fmt.Errorf("unknown writer kind: %s", cfg.Writer)
	}
}

// convertHeaders converts AMQP headers to a map[string]interface{}
func convertHeaders(amqpHeaders amqp.Table) map[string]interface{} {
	if amqpHeaders == nil {
		return nil
	}

	headers := make(map[string]interface{})
	for k, v := range amqpHeaders {
		headers[k] = v
	}
	return headers
}

// writeMessageCommon handles the message creation and serialization
func writeMessageCommon(msg amqp.Delivery, prettyPrint bool) ([]byte, error) {
	// Create a new Message struct with AMQP delivery details
	message := model.Message{
		Headers:    convertHeaders(msg.Headers),
		Exchange:   msg.Exchange,
		RoutingKey: msg.RoutingKey,
		Body:       msg.Body,
	}

	// Prepare output based on pretty print config
	var output []byte
	var err error
	if prettyPrint {
		output, err = json.MarshalIndent(message, "", "  ")
	} else {
		output, err = json.Marshal(message)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	// Always append a newline
	output = append(output, '\n')

	return output, nil
}
