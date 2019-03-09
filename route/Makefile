all: build-docker
build:
	CGO_ENABLED=0 go build -o route -mod=vendor .

shell:
	docker run -it -v $(shell pwd):/build -w /build golang:1.11

build-docker:
	docker run -it -v $(shell pwd):/build -w /build golang:1.11 make build

vendor:
	go mod tidy
	go mod vendor

vendor-docker:
	docker run -it -v $(shell pwd):/build -w /build golang:1.11 make vendor
