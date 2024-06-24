FROM golang:1.22-alpine as builder
WORKDIR /workspace
RUN apk update && apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webhook main.go
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /workspace/webhook .
COPY --from=builder /workspace/config ./config
EXPOSE 443
CMD ["./webhook", "-config=config/config.yaml", "-certFile=/etc/webhook/certs/tls.crt", "-keyFile=/etc/webhook/certs/tls.key", "-caFile=/etc/webhook/certs/ca.crt"]
