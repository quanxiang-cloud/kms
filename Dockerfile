FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o kms -mod=vendor -ldflags='-s -w' -installsuffix cgo cmd/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /kms
COPY --from=builder ./build/kms ./cmd/

EXPOSE 80

ENTRYPOINT ["./cmd/kms", "-config=/configs/config.yml"]