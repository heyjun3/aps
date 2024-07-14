package main

import (
	"flag"
	"log"
	"os"

	"crawler/db"
)

func main() {
	var (
		target string
	)
	flag.StringVar(&target, "t", "prd", "expect environment")
	flag.Parse()

	switch {
	case target == "prd":
		db.RunMigrate(os.Getenv("DB_DSN"))
	case target == "test":
		db.RunMigrate(os.Getenv("TEST_DSN"))
	default:
		log.Fatalf("argument error: t=%s", target)
	}
}
