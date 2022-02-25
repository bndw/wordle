# WordleSSH

Wordle over SSH

## Run locally

Run the server on :2222

```
make build run
```

From another terminal, connect:

```
make ssh
```

To generate a key.pem cert, use this:
https://go.dev/src/crypto/tls/generate_cert.go

## How to play

```
ssh wordle.bdw.to
```
