package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// EventStatusChanged represents the event structure consumed by the service
type EventStatusChanged struct {
	RegistrationID string `json:"registration_id"`
	Status         string `json:"status"`
	Timestamp      string `json:"timestamp"`
}

func main() {
	// Command-line flags for flexibility
	broker := flag.String("broker", "localhost:19093", "Kafka broker address")
	topic := flag.String("topic", "event.status.changed", "Kafka topic")
	registrationID := flag.String("id", "REG-001", "Registration ID")
	status := flag.String("status", "confirmed", "Event status (confirmed, cancelled, pending, etc.)")
	count := flag.Int("count", 1, "Number of messages to send")
	flag.Parse()

	// Create Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{*broker},
		Topic:    *topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	ctx := context.Background()

	fmt.Printf("🚀 Kafka Producer Started\n")
	fmt.Printf("📡 Broker: %s\n", *broker)
	fmt.Printf("📬 Topic: %s\n", *topic)
	fmt.Printf("📨 Sending %d message(s)...\n\n", *count)

	// Send messages
	for i := 0; i < *count; i++ {
		// Create event
		event := EventStatusChanged{
			RegistrationID: fmt.Sprintf("%s-%d", *registrationID, i+1),
			Status:         *status,
			Timestamp:      time.Now().Format(time.RFC3339),
		}

		// Marshal to JSON
		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("❌ Failed to marshal event: %v", err)
		}

		// Send message
		msg := kafka.Message{
			Key:   []byte(event.RegistrationID),
			Value: eventBytes,
		}

		err = writer.WriteMessages(ctx, msg)
		if err != nil {
			log.Printf("❌ Failed to send message %d: %v", i+1, err)
			continue
		}

		fmt.Printf("✅ Message %d sent successfully:\n", i+1)
		fmt.Printf("   Registration ID: %s\n", event.RegistrationID)
		fmt.Printf("   Status: %s\n", event.Status)
		fmt.Printf("   Timestamp: %s\n", event.Timestamp)
		fmt.Printf("   Payload: %s\n\n", string(eventBytes))

		// Small delay between messages if sending multiple
		if *count > 1 && i < *count-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("🎉 Successfully sent %d message(s) to topic '%s'\n", *count, *topic)
}
