# Build Stage
FROM golang:1.16-alpine as builder
RUN apk --no-cache add git build-base

WORKDIR /go/src/github.com/bndw/wordle
COPY go.* ./
RUN go mod tidy
COPY . .

RUN go build -o /bin/app .

# Execution Stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/app .

EXPOSE 22
CMD ["./app"]
