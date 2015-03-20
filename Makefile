CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG=${TAG:-latest}

all: deps build

deps:
	@godep restore

clean:
	@rm -rf Godeps/_workspace machine-server

build:
	@godep go build -a -tags 'netgo' -ldflags '-w -linkmode external -extldflags -static' .

image: build
	@echo Building image $(TAG)
	@docker build -t ehazlett/machine-server:$(TAG) .

release: deps build image
	@docker push ehazlett/machine-server:$(TAG)

test: clean 
	@godep go test -v ./...

.PHONY: all deps build clean image test release
