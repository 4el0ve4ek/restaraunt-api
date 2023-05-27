FROM golang:1.20 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app

COPY . .
RUN go mod download

WORKDIR /usr/src/app/cmd/auth
RUN go build -o /usr/local/bin/auth -buildvcs=false

FROM alpine:latest

COPY --from=builder /usr/local/bin/auth ./
COPY --from=builder /usr/src/app/config.yaml ./config.yaml
CMD ["./auth"]
EXPOSE 8080