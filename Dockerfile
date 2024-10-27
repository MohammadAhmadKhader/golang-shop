FROM golang:1.23.2

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/

EXPOSE 8080

ENTRYPOINT [ "/api" ]