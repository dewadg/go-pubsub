package main

import (
	"context"
	"log"
)

type pubsubImpl struct {
	queue       chan pubsubEvent
	subscribers []pubsubHandler
}

func (pubsub *pubsubImpl) Publish(ctx context.Context, payload string) {
	go func() {
		pubsub.queue <- pubsubEvent{
			Context: ctx,
			Payload: payload,
		}
	}()
}

func (pubsub *pubsubImpl) RegisterSubscriber(handler pubsubHandler) {
	pubsub.subscribers = append(pubsub.subscribers, handler)
}

func (pubsub *pubsubImpl) Run() {
	go func() {
		for {
			select {
			case event := <-pubsub.queue:
				for _, subscriber := range pubsub.subscribers {
					if err := subscriber(event.Context, event.Payload); err != nil {
						log.Println(err.Error())
					}
				}
				break
			}
		}
	}()
}

func (pubsub *pubsubImpl) Stop() {
	close(pubsub.queue)
}
