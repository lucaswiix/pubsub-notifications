FROM golang:1 as builder
WORKDIR /go/src/github.com/lucaswiix/meli/notifications

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:latest  
RUN apk --no-cache upgrade
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /go/src/github.com/lucaswiix/meli/notifications/api .
COPY --from=builder /go/src/github.com/lucaswiix/meli/notifications/api.yaml .
CMD ["./api"]%