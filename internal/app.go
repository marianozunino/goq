package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/filewriter"
	"github.com/marianozunino/goq/internal/rmq"
)

// Dump contains the main application logic
func Dump(cfg *config.Config) error {
	if err := printConfig(cfg); err != nil {
		return err
	}

	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		return handleError("failed to create consumer", err)
	}
	defer consumer.Close()

	writer, err := filewriter.NewWriter(cfg)
	if err != nil {
		return handleError("failed to create file writer", err)
	}
	defer writer.Close()

	msgCount, err := consumer.GetQueueInfo()
	if err != nil {
		return handleError("failed to get queue info", err)
	}
	color.Cyan("Number of messages in the queue: %d", msgCount)

	msgs, err := consumer.ConsumeMessages()
	if err != nil {
		return handleError("failed to consume messages", err)
	}

	if cfg.StopAfterConsume {
		color.Yellow("Stopping after consuming all messages")
	}

	log.Println("Waiting for messages. To exit press CTRL+C")
	return processMessages(msgs, writer, msgCount, cfg.StopAfterConsume)
}

// Monitor monitors messages on a temporary queue
func Monitor(cfg *config.Config) error {
	if err := printConfig(cfg); err != nil {
		return err
	}

	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		return handleError("failed to create consumer", err)
	}
	defer consumer.Close()

	tempQueue, err := consumer.DeclareTemporaryQueue()
	if err != nil {
		return handleError("failed to declare temporary queue", err)
	}

	if err := bindQueueToRoutingKeys(consumer, tempQueue.Name, cfg.RoutingKeys); err != nil {
		return err
	}

	color.Cyan("Monitoring messages on temporary queue: %s", tempQueue.Name)

	msgs, err := consumer.ConsumeMessagesFromQueue(tempQueue.Name)
	if err != nil {
		return handleError("failed to consume messages", err)
	}

	writer, err := filewriter.NewWriter(cfg)
	if err != nil {
		return handleError("failed to create file writer", err)
	}
	defer writer.Close()

	return logReceivedMessages(msgs, writer, cfg.PrettyPrint)
}

// Helper function to print the configuration
func printConfig(cfg *config.Config) error {
	configJSON, err := cfg.PrintConfig()
	if err != nil {
		return fmt.Errorf("failed to print config: %v", err)
	}
	color.Green("Configuration used:")
	fmt.Println(configJSON)
	return nil
}

// Helper function to handle errors
func handleError(msg string, err error) error {
	return fmt.Errorf("%s: %v", msg, err)
}

// Helper function to bind the temporary queue to routing keys
func bindQueueToRoutingKeys(consumer *rmq.Consumer, queueName string, routingKeys []string) error {
	color.Green("Binding queue %q to routing keys", queueName)
	for _, routingKey := range routingKeys {
		routingKeyTrimmed := strings.TrimSpace(routingKey)
		if err := consumer.BindQueue(queueName, routingKeyTrimmed); err != nil {
			return handleError(fmt.Sprintf("failed to bind queue %q to routing key %q", queueName, routingKeyTrimmed), err)
		}
		color.Yellow("Bound routing key %q", routingKeyTrimmed)
	}
	return nil
}

// Helper function to process messages
func processMessages(msgs <-chan rmq.Message, writer *filewriter.Writer, msgCount int, stopAfterConsume bool) error {
	if msgCount == 0 && stopAfterConsume {
		color.Yellow("No messages to consume. Exiting.")
		return nil
	}

	consumedCount := 0
	blue := color.New(color.FgBlue)

	for msg := range msgs {
		if err := writer.WriteMessage(string(msg.Body)); err != nil {
			log.Printf("Failed to write message: %v", err)
		}
		consumedCount++
		blue.Printf("\rMessages dumped: %d/%d", consumedCount, msgCount)
		os.Stdout.Sync() // Flush the output

		if stopAfterConsume && consumedCount >= msgCount {
			break
		}
	}
	fmt.Println()
	color.Green("All messages have been consumed. Exiting.")
	return nil
}

// Helper function to log received messages
func logReceivedMessages(msgs <-chan rmq.Message, writer *filewriter.Writer, prettyPrint bool) error {
	blue := color.New(color.FgBlue)

	for msg := range msgs {
		if err := writer.WriteMessage(string(msg.Body)); err != nil {
			log.Printf("Failed to write message: %v", err)
		}
		if prettyPrint {
			blue.Printf("%s\n", prettyPrintJson(msg.Body))
		} else {
			blue.Printf("%s\n", string(msg.Body))
		}
		msg.Ack(true)
	}

	color.Green("Consumer disconnected. Exiting monitoring.")
	return nil
}

func prettyPrintJson(b []byte) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, b, "", "  ")
	if error != nil {
		log.Println("Failed to parse JSON. Error:", error)
	}

	return prettyJSON.String()
}
