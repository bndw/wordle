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

## TODO

- [x] Use official word list
- [x] Persist data
- [ ] Save partial game state (write game on each guess)
- [x] AddHostKey to avoid known_hosts issues: https://pkg.go.dev/github.com/gliderlabs/ssh#Server.AddHostKey
- [x] Limit players once per day
- [ ] Display game stats over http (html)

## How to play

```
ssh wordle.bdw.to
```
