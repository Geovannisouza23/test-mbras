FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server ./server
EXPOSE 8080
ENTRYPOINT ["./server"]
