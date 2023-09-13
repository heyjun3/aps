package server

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"

	greetv1 "crawler/server/gen/greet/v1"
)

type GreetServer struct {}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Msg.Name)
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
