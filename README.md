# Wordle over SSH

Play Wordle over SSH!

```
ssh wordle.bdw.to
```

## Run locally

To run the server locally on :2222, first generate a PEM file `key.pem` in this directory and then run:

```
make build run
```

From another terminal, connect:

```
make ssh
```
