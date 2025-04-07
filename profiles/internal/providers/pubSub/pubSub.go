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
	fmt.Println("PublishToService called")
	messageData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("PublishToService json marshal error")
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	topic := ps.client.Topic(serviceName)
	fmt.Println("PublishToService topic: ", topic)
	msg := &pubsub.Message{
		Data: messageData,
	}

	result := topic.Publish(ctx, msg)
	fmt.Println("PublishToService published")
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("PublishToService errored", err)
		return fmt.Errorf("failed to publish message: %v", err)
	}

	fmt.Printf("PublishToService Message published with ID: %s", id)
	return nil
}
