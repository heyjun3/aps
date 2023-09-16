package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"

	greetv1 "crawler/server/gen/greet/v1"
	"crawler/server/gen/greet/v1/greetv1connect"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
	log.Println(res.Header())
}
