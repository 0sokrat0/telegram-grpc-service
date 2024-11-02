FROM golang:1.22 AS builder


WORKDIR /app
COPY . .
RUN go build -o gateway ./cmd/gateway

FROM gcr.io/distroless/base
COPY --from=builder /app/gateway /gateway
ENTRYPOINT ["/gateway"]
