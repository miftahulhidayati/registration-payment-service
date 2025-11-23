package kafka

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
    "example.com/registration-payment-service/internal/repository"
)

type EventStatusChanged struct {
    RegistrationID string `json:"registration_id"`
    Status         string `json:"status"`
    Timestamp      string `json:"timestamp"`
}

// StartConsumer starts a Kafka consumer for the event.status.changed topic.
// It reads messages, unmarshals them into EventStatusChanged, and updates the registration status in the database.
func StartConsumer(brokers string, topic string, groupID string, repo *repository.Postgres) error {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{brokers},
        Topic:     topic,
        GroupID:   groupID,
        MinBytes:  10e3, // 10KB
        MaxBytes:  10e6, // 10MB
        CommitInterval: time.Second,
    })
    defer r.Close()

    ctx := context.Background()
    for {
        m, err := r.ReadMessage(ctx)
        if err != nil {
            log.Printf("error reading kafka message: %v", err)
            return err
        }
        var evt EventStatusChanged
        if err := json.Unmarshal(m.Value, &evt); err != nil {
            log.Printf("failed to unmarshal event.status.changed: %v", err)
            continue
        }
        // Example: update registration status in DB (implementation placeholder)
        if err := updateRegistrationStatus(ctx, repo, evt); err != nil {
            log.Printf("failed to update registration status: %v", err)
        } else {
            log.Printf("processed event.status.changed for registration %s, status %s", evt.RegistrationID, evt.Status)
        }
    }
}

func updateRegistrationStatus(ctx context.Context, repo *repository.Postgres, evt EventStatusChanged) error {
    // Placeholder implementation – you should replace this with actual DB logic.
    // For example, you might have a method repo.UpdateRegistrationStatus(id, status).
    // Here we just log the intention.
    log.Printf("[DB] Update registration %s to status %s (timestamp %s)", evt.RegistrationID, evt.Status, evt.Timestamp)
    return nil
}


