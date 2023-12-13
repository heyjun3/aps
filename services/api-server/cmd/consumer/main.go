package main

import (
	"api-server/mq"
)

func main() {
	mq.Consume(mq.Client{}.Exec, "chart")
}
