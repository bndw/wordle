REPO ?= bndw/wordle
GITSHA=$(shell git rev-parse --short HEAD)
TAG_COMMIT=$(REPO):$(GITSHA)
TAG_LATEST=$(REPO):latest

.PHONY: all
all: dev

.PHONY: dev
dev: build run

.PHONY: build
build:
	@docker build -t $(REPO) .

build-bin:
	go build -o ./bin/wordle .

.PHONY: publish
publish:
	docker push $(TAG_LATEST)
	@docker tag $(TAG_LATEST) $(TAG_COMMIT)
	docker push $(TAG_COMMIT)

.PHONY: run
run: wordle.db
	@echo 'Listening on localhost:2222'
	@docker run \
		-t \
		--rm \
		-p 2222:22 \
		-v $(shell pwd)/wordle.db:/root/wordle.db \
		$(REPO)

wordle.db:
	@touch $@

buildlinux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/atm .

alaska: buildlinux
	scp ./bin/atm alaska:~/

gorun:
	./bin/game

ssh:
	ssh -p 2222 -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa localhost


test:
	go test -v ./...
