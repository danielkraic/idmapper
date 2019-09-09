BINARY_NAME := idmapperapp
DOCKER_REPOSITORY := danielkraic
DOCKER_IMAGE := idmapper

VERSION := $(shell git describe --tags 2> /dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD := $(shell date +"%FT%T%z")

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Commit=$(COMMIT) -X=main.Build=$(BUILD)"

all: test race cover build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $(BINARY_NAME) -v

clean:
	go clean
	rm -f $(BINARY_NAME)

test: fmt vet lint errcheck
	go test -v ./...

race:
	go test -race ./...

cover:
	go test -cover ./...

coverprofile:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out
	
errcheck:
	go list ./... | grep -v "vendor\|errcheck" | xargs errcheck

vet:
	go list ./... | grep -v vendor | xargs go vet

lint:
	go list ./... | grep -v vendor | xargs golint

fmt:
	go fmt ./...

docker:
	docker build -t $(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(VERSION) .