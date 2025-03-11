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

up-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d $(ARGS)

up-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d $(ARGS)

up-test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up -d $(ARGS)

down-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down $(ARGS)

down-dev:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down $(ARGS)

down-test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml down $(ARGS)

swarm-up:
	docker stack deploy -c docker-compose.yml -c docker-compose.prod.yml golang-shop-app

swarm-down:
	docker stack rm golang-shop-app

test-nginx:
	bash -c "./curlreqs.sh"