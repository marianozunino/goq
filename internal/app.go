package app

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/filewriter"
	"github.com/marianozunino/goq/internal/rmq"
)

// Run contains the main application logic
func Run(cfg *config.Config) error {
	configJSON, err := cfg.PrettyPrint()
	if err != nil {
		return fmt.Errorf("failed to print config: %v", err)
	}
	color.Green("Configuration used:")
	fmt.Println(configJSON)

	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	writer, err := filewriter.NewWriter(cfg)
	if err != nil {
		return fmt.Errorf("failed to create file writer: %v", err)
	}
	defer writer.Close()

	msgCount, err := consumer.GetQueueInfo()
	if err != nil {
		return fmt.Errorf("failed to get queue info: %v", err)
	}
	color.Cyan("Number of messages in the queue: %d", msgCount)

	msgs, err := consumer.ConsumeMessages()
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	if cfg.StopAfterConsume {
		color.Yellow("Stopping after consuming all messages")
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	consumedCount := 0
	blue := color.New(color.FgBlue)
	for msg := range msgs {
		err := writer.WriteMessage(string(msg.Body))
		if err != nil {
			log.Printf("Failed to write message: %v", err)
		}

		consumedCount++
		blue.Printf("\rMessages dumped: %d/%d", consumedCount, msgCount)
		os.Stdout.Sync() // Flush the output

		if cfg.StopAfterConsume && consumedCount >= msgCount {
			break
		}
	}
	fmt.Println()

	color.Green("All messages have been consumed. Exiting.")
	return nil
}
