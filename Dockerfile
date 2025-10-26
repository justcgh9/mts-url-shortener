FROM golang:1.24 AS builder
WORKDIR /app

ARG SERVER_PORT
ARG CONFIG_PATH

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY ${CONFIG_PATH} config/prod.yml

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/url-shortener ./cmd/url-shortener

FROM alpine:latest
WORKDIR /app

ARG SERVER_PORT

COPY --from=builder /app/bin/url-shortener ./bin/url-shortener
COPY --from=builder /app/config/prod.yml ./config/prod.yml

EXPOSE ${SERVER_PORT}

CMD ["./bin/url-shortener", "--config=./config/prod.yml"]