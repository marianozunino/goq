package app

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/exporter"
	"github.com/marianozunino/goq/internal/rmq"
	"github.com/streadway/amqp"
)

// MessageProcessor handles the core logic of processing messages
type MessageProcessor struct {
	config   *config.Config
	consumer *rmq.Consumer
	exporter exporter.Exporter
	msgCount int
}

func NewMessageProcessor(cfg *config.Config) (*MessageProcessor, error) {
	// Print configuration
	configJSON, err := cfg.PrintConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to print config: %v", err)
	}

	color.Green("Configuration used:")
	fmt.Println(configJSON)

	// Create consumer
	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to create file exporter: %v", err)
	}

	// Get queue message count (only for non-empty queue name)
	var msgCount int
	if cfg.Queue != "" {
		msgCount, err = consumer.GetQueueInfo()
		if err != nil {
			consumer.Close()
			exp.Close()
			return nil, fmt.Errorf("failed to get queue info: %v", err)
		}
		// Print queue information only if queue exists
		color.Cyan("Number of messages in the queue: %d", msgCount)
	}

	return &MessageProcessor{
		config:   cfg,
		consumer: consumer,
		exporter: exp,
		msgCount: msgCount,
	}, nil
}

// Dump processes messages from the main queue
func (mp *MessageProcessor) Dump() error {
	defer mp.consumer.Close()
	defer mp.exporter.Close()

	// Consume messages
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
	defer mp.consumer.Close()
	defer mp.exporter.Close()

	// Consume messages from temporary queue
	msgs, err := mp.consumer.Consume()
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	return mp.processMessages(msgs)
}

// processMessages handles the core message processing logic
func (mp *MessageProcessor) processMessages(msgs <-chan amqp.Delivery) error {
	if mp.msgCount == 0 && mp.config.StopAfterConsume {
		color.Yellow("No messages to consume. Exiting.")
		return nil
	}

	blue := color.New(color.FgBlue)
	consumedCount := 0

	for msg := range msgs {
		// Write message to exporter
		if err := mp.exporter.WriteMessage(msg); err != nil {
			log.Printf("Failed to write message: %v", err)
			continue
		}

		consumedCount++

		// Handle different output formats
		switch mp.config.Writer {
		case config.FileWriterKind:
			blue.Printf("\rMessages processed: %d/%d", consumedCount, mp.msgCount)
		case config.ConsoleExporterKind:
			blue.Println("*****")
		}

		// Stop after consuming all messages if configured
		if mp.config.StopAfterConsume && consumedCount >= mp.msgCount {
			break
		}
	}

	fmt.Println()
	color.Green("Message processing complete.")
	return nil
}
