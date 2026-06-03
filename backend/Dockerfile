FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pbx-backend main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pbx-backend .

EXPOSE 8081
CMD ["./pbx-backend"]
