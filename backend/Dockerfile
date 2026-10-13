FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o pbx-backend .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/pbx-backend .

EXPOSE 8081
CMD ["./pbx-backend"]
