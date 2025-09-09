package exporter

import (
	"encoding/json"
	"fmt"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/model"
	"github.com/wagslane/go-rabbitmq"
)

// ExporterError represents different types of exporter errors
type ExporterError struct {
	Type string
	Err  error
}

func (e *ExporterError) Error() string {
	return fmt.Sprintf("%s error: %v", e.Type, e.Err)
}

func (e *ExporterError) Unwrap() error {
	return e.Err
}

// Error types
const (
	ErrorTypeSerialization = "serialization"
	ErrorTypeFileIO        = "file_io"
	ErrorTypeConsoleIO     = "console_io"
	ErrorTypeConfiguration = "configuration"
)

type Exporter interface {
	WriteMessage(msg rabbitmq.Delivery) error
	Close() error
}

// ExporterFactory defines the interface for creating exporters
type ExporterFactory interface {
	CreateExporter(cfg *config.Config) (Exporter, error)
	GetType() string
}

// DefaultExporterFactory implements the factory pattern
type DefaultExporterFactory struct{}

func (f *DefaultExporterFactory) GetType() string {
	return "default"
}

func (f *DefaultExporterFactory) CreateExporter(cfg *config.Config) (Exporter, error) {
	switch cfg.Writer {
	case config.ConsoleExporterKind:
		return NewConsoleExporter(cfg)
	case config.FileWriterKind:
		return NewFileWriter(cfg)
	default:
		return nil, &ExporterError{
			Type: ErrorTypeConfiguration,
			Err:  fmt.Errorf("unknown writer kind: %s", cfg.Writer),
		}
	}
}

// Registry for exporter factories
var exporterFactories = map[string]ExporterFactory{
	"default": &DefaultExporterFactory{},
}

// RegisterExporterFactory allows adding new exporter types
func RegisterExporterFactory(name string, factory ExporterFactory) {
	exporterFactories[name] = factory
}

// NewExporter creates an exporter using the default factory
func NewExporter(cfg *config.Config) (Exporter, error) {
	factory := exporterFactories["default"]
	return factory.CreateExporter(cfg)
}

// convertHeaders converts AMQP headers to a map[string]interface{}
func convertHeaders(amqpHeaders rabbitmq.Table) map[string]interface{} {
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
func writeMessageCommon(msg rabbitmq.Delivery, prettyPrint bool) ([]byte, error) {
	message := model.Message{
		Headers:    convertHeaders(rabbitmq.Table(msg.Headers)),
		Exchange:   msg.Exchange,
		RoutingKey: msg.RoutingKey,
		Body:       msg.Body,
	}

	var output []byte
	var err error
	if prettyPrint {
		output, err = json.MarshalIndent(message, "", "  ")
	} else {
		output, err = json.Marshal(message)
	}

	if err != nil {
		return nil, &ExporterError{
			Type: ErrorTypeSerialization,
			Err:  fmt.Errorf("failed to marshal message: %v", err),
		}
	}

	output = append(output, '\n')

	return output, nil
}
