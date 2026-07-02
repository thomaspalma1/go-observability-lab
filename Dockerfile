FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/api ./cmd/api

FROM alpine:latest
COPY --from=builder /bin/api /bin/api
EXPOSE 8082
ENTRYPOINT ["/bin/api"]
