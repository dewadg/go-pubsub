package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	var pubsub pubsubContract

	pubsub = &pubsubImpl{
		queue:       make(chan pubsubEvent),
		subscribers: make([]pubsubHandler, 0),
	}

	// Run
	pubsub.Run()

	// Register some subscribers
	pubsub.RegisterSubscriber(func(ctx context.Context, s string) error {
		fmt.Println(s)

		return nil
	})
	pubsub.RegisterSubscriber(func(ctx context.Context, s string) error {
		fmt.Println(strings.ToUpper(s))

		return nil
	})
	pubsub.RegisterSubscriber(func(ctx context.Context, s string) error {
		fmt.Println(strings.ToLower(s))

		return nil
	})
	pubsub.RegisterSubscriber(func(ctx context.Context, s string) error {
		return fmt.Errorf("Error with payload %s", s)
	})

	// Serve a HTTP handler
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		payload := request.URL.Query().Get("payload")
		if payload == "" {
			http.Error(writer, "Query `payload` is required", http.StatusBadRequest)
			return
		}

		pubsub.Publish(request.Context(), payload)
		fmt.Fprint(writer, fmt.Sprintf("Payload `%s` received", payload))
	})

	fmt.Println("[*] Running on port 8000")
	http.ListenAndServe("localhost:8000", nil)
}
