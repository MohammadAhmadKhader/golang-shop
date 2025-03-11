FROM golang:1.23.2 as builder

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api/

FROM alpine:latest as runtime
WORKDIR /root
COPY --from=builder /app/api /root/api
RUN apk --no-cache add bash
EXPOSE 8080

CMD [ "/root/api" ]

# FROM golang:1.23.2

# WORKDIR /app
# COPY go.mod go.sum ./
# COPY . .
# RUN go mod download

# RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/

# EXPOSE 8080

# ENTRYPOINT [ "/api" ]