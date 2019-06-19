SHELL = bash
PROJECT := oasis

COMMIT := $(shell git rev-parse --short HEAD)
DOCKERIMAGE ?= $(shell zutano docker image --name=$(PROJECT))

all: binaries

clean:
	rm -Rf bin

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64 windows/amd64" \
		-ldflags="-X main.projectVersion=${VERSION} -X main.projectBuild=${COMMIT}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		./...

.PHONY: test
test:
	mkdir -p bin/test
	go test -coverprofile=bin/test/coverage.out -v ./... | tee bin/test/test-output.txt ; exit "$${PIPESTATUS[0]}"
	cat bin/test/test-output.txt | go-junit-report > bin/test/unit-tests.xml
	go tool cover -html=bin/test/coverage.out -o bin/test/coverage.html

bootstrap:
	go get github.com/arangodb-managed/zutano
	go get github.com/mitchellh/gox
	go get github.com/jstemmer/go-junit-report

docker:
	docker build \
		--build-arg=GOARCH=amd64 \
		-t $(DOCKERIMAGE) .

docker-push:
	docker push $(DOCKERIMAGE)


.PHONY: update-modules
update-modules:
	rm -f go.mod go.sum 
	go mod init
	go mod edit \
		-replace github.com/golang/lint=golang.org/x/lint@v0.0.0-20181026193005-c67002cb31c3
	go get -u \
		github.com/arangodb-managed/apis@v0.12.1
	go mod tidy
