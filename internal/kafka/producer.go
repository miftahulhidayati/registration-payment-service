package kafka

import (
    "context"
    "encoding/json"
    "time"

    kgo "github.com/segmentio/kafka-go"
)

type Producer struct {
    writer *kgo.Writer
}

func NewProducer(brokers string) (*Producer, error) {
    if brokers == "" {
        return nil, nil
    }
    w := &kgo.Writer{
        Addr:     kgo.TCP(brokers),
        Balancer: &kgo.LeastBytes{},
    }
    return &Producer{writer: w}, nil
}

func (p *Producer) Close() error {
    if p == nil || p.writer == nil {
        return nil
    }
    return p.writer.Close()
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, data any) error {
    if p == nil || p.writer == nil || topic == "" {
        return nil
    }
    b, err := json.Marshal(data)
    if err != nil {
        return err
    }
    msg := kgo.Message{
        Topic: topic,
        Key:   []byte(key),
        Value: b,
        Time:  time.Now(),
    }
    return p.writer.WriteMessages(ctx, msg)
}


