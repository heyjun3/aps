
fmt:
	go fmt ./...

migrate-run:
	go run cmd/migrate/main.go

migrate-rollback:
	go run cmd/migrate/main.go -r

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(FILE)
