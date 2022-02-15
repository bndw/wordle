# Build Stage
FROM golang:1.13-alpine as builder
RUN apk --no-cache add git

WORKDIR /go/src/github.com/bndw/wordle
COPY . .

RUN go build -o /bin/app .

# Execution Stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/app .

EXPOSE 22
CMD ["./app"]
