FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o risk-engine ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/risk-engine /app/risk-engine
COPY --from=builder /app/model /app/model
EXPOSE 8081
CMD ["/app/risk-engine"]
