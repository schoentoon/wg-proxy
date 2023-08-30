FROM docker.io/library/golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/wg-proxy ./cmd/wg-proxy/...

FROM gcr.io/distroless/base:latest

COPY --from=builder /bin/wg-proxy /bin/wg-proxy

CMD ["/bin/wg-proxy"]
