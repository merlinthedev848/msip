FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -mod=vendor -o pbx-backend .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pbx-backend .

EXPOSE 8081
CMD ["./pbx-backend"]
