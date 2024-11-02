FROM golang:1.22 AS builder



WORKDIR /app
COPY . .
RUN go build -o grpc_server ./cmd/server

FROM gcr.io/distroless/base
COPY --from=builder /app/grpc_server /grpc_server
ENTRYPOINT ["/grpc_server"]
