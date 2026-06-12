# ---- build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

# ---- run stage ----
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]
