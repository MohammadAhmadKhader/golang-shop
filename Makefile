seed:
	@go run ./cmd/seed/main.go

backup:
	@go run ./cmd/backup/ $(ARGS)

build:
	@go build -o ./bin/golang_shop ./cmd/api/main.go

run: build
	@./bin/golang_shop

test:
	@go test -v ./... -count=1