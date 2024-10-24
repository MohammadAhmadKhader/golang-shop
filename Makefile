.PHONY: backup

run:
	@go run ./cmd/api/main.go

seed:
	@go run ./cmd/seed/main.go

backup:
	@go run ./cmd/backup/ $(ARGS)
