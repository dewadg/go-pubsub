package main

import "context"

type pubsubEvent struct {
	Context context.Context
	Payload string
}

type pubsubHandler func(context.Context, string) error

type pubsubContract interface {
	Publish(context.Context, string)

	RegisterSubscriber(pubsubHandler)

	Run()

	Stop()
}
