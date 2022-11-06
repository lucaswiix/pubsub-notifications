FROM golang:1 as builder
WORKDIR /go/src/github.com/lucaswiix/notifications-tracking-app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bff ./cmd/app

FROM alpine:latest  
RUN apk --no-cache upgrade
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /go/src/github.com/lucaswiix/notifications-tracking-app/bff .
COPY --from=builder /go/src/github.com/lucaswiix/notifications-tracking-app/config.yaml .
COPY --from=builder /go/src/github.com/lucaswiix/notifications-tracking-app/website .
CMD ["./bff"]%