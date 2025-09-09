package app

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/exporter"
	"github.com/marianozunino/goq/internal/rmq"
)

// MessageProcessor handles the core logic of processing messages
type MessageProcessor struct {
	config   *config.Config
	consumer *rmq.Consumer
	exporter exporter.Exporter
}

// NewMessageProcessor creates a new MessageProcessor
func NewMessageProcessor(cfg *config.Config) (*MessageProcessor, error) {
	// Print configuration
	configtable := cfg.PrintConfig()

	color.Green("Configuration used:")
	fmt.Println(configtable)

	// Create consumer
	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create file exporter: %v", err)
	}

	return &MessageProcessor{
		config:   cfg,
		consumer: consumer,
		exporter: exp,
	}, nil
}

// Dump processes messages from the main queue
func (mp *MessageProcessor) Dump() error {
	defer mp.exporter.Close()

	msgs, err := mp.consumer.Consume()
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	if mp.config.StopAfterConsume {
		color.Yellow("Stopping after consuming all messages")
	}
	log.Println("Waiting for messages. To exit press CTRL+C")

	return mp.processMessages(msgs)
}

// Monitor creates a temporary queue and processes messages
func (mp *MessageProcessor) Monitor() error {
	defer mp.exporter.Close()

	// Consume messages from temporary queue
	msgs, err := mp.consumer.Consume()
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	return mp.endlessConsume(msgs)
}

func (mp *MessageProcessor) processMessages(status <-chan rmq.ConsumerStatus) error {
	blue := color.New(color.FgBlue)
	for s := range status {
		// when message is null is because the message was filtered
		if s.Message != nil {
			if err := mp.exporter.WriteMessage(*s.Message); err != nil {
				log.Printf("Failed to write message: %v", err)
				continue
			}

			switch mp.config.Writer {
			case config.FileWriterKind:
				blue.Printf("\rMessages processed: %d", s.ConsumedMessages)
			case config.ConsoleExporterKind:
				blue.Println("*****")
			}
		}

		if s.Complete {
			fmt.Println()
			color.Green("Message processing complete.")
			return nil
		}
	}
	return nil
}

func (mp *MessageProcessor) endlessConsume(status <-chan rmq.ConsumerStatus) error {
	blue := color.New(color.FgBlue)
	for s := range status {
		// when message is null is because the message was filtered
		if s.Message != nil {
			if err := mp.exporter.WriteMessage(*s.Message); err != nil {
				log.Printf("Failed to write message: %v", err)
				continue
			}

			switch mp.config.Writer {
			case config.FileWriterKind:
				blue.Printf("\rMessages processed: %d", s.ConsumedMessages)
			case config.ConsoleExporterKind:
				blue.Println("*****")
			}
		}
	}
	return nil
}
