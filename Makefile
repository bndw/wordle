build:
	go build -o ./bin/wordle .

run: wordle.db
	./bin/wordle -cert key.pem -port 2222

wordle.db:
	@touch $@

ssh:
	ssh -p 2222 -i ~/.ssh/id_rsa localhost

test:
	go test -v ./...
