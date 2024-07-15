package main

import (
	"flag"
	"log"
	"os"

	"crawler/db"
)

func main() {
	var (
		target     string
		isRollback bool
	)
	flag.StringVar(&target, "t", "prd", "expect environment")
	flag.BoolVar(&isRollback, "r", false, "expect rollback")
	flag.Parse()

	steps := 1
	if isRollback {
		steps = -1
	}

	switch {
	case target == "prd":
		db.RunMigrate(os.Getenv("DB_DSN"), steps)
	case target == "test":
		db.RunMigrate(os.Getenv("TEST_DSN"), steps)
	default:
		log.Fatalf("argument error: t=%s", target)
	}
}
