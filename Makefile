GOCMD=go
BINARY_NAME=simple-uptime-backend

.PHONY: all build vendor

vendor:
	$(GOCMD) mod vendor

build:
	mkdir -p out/bin
	GO111MODULE=on $(GOCMD) build -mod vendor -o out/bin/$(BINARY_NAME) .
