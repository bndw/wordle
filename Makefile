build:
	go build -o ./bin/wordle .

run: wordle.db
	./bin/wordle -key key.pem -port 2222

wordle.db:
	@touch $@

ssh:
	ssh -p 2222 localhost

test:
	go test -v ./...
