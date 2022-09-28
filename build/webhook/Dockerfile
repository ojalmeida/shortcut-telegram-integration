FROM alpine:latest AS builder

RUN apk add go ca-certificates

RUN mkdir -p /usr/local/src

WORKDIR /usr/local/src

COPY src .

RUN CGO_ENABLED=0 GOOS=linux go build -o out -a -installsuffix cgo cmd/main.go

FROM alpine:latest

COPY --from=builder /usr/local/src/out /usr/local/bin/webhook

ENTRYPOINT /usr/local/bin/webhook