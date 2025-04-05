package PubSub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

type PubSub struct {
	client *pubsub.Client
}

var pubSub *PubSub
var once sync.Once

func Init(ctx context.Context, projectID string) {
	once.Do(func() {
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			panic(err)
		}
		pubSub = &PubSub{client: client}
	})
}

func GetClient() *PubSub {
	if pubSub == nil {
		panic("PubSub client not initialized, call Init(...) first")
	}
	return pubSub
}

func (ps *PubSub) PublishToService(ctx context.Context, serviceName string, message PublishMessageType) error {
	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	topic := ps.client.Topic(serviceName)
	msg := &pubsub.Message{
		Data: messageData,
	}
	result := topic.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("Failed to get message", err)
		return fmt.Errorf("failed to publish message: %v", err)
	}

	fmt.Printf("Message published with ID: %s", id)
	return nil
}
