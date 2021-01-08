SHELL = bash
PROJECT := oasisctl

COMMIT := $(shell zutano repo build)
VERSION := $(shell zutano repo version)
DOCKERIMAGE ?= $(shell zutano docker image --name=$(PROJECT))
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GOEXE := $(shell go env GOEXE)

all: binaries

clean:
	rm -Rf bin assets

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64 windows/amd64" \
		-ldflags="-X main.projectVersion=${VERSION} -X main.projectBuild=${COMMIT}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		./...
	mkdir -p assets
	zip -r assets/oasisctl.zip bin/*

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

publish-oasis-tools:
	GITHUB_USERNAME=$(CIRCLE_PROJECT_USERNAME) COMMIT=$(COMMIT) VERSION=$(VERSION) ./scripts/publish-oasis-tools.sh

update-apis-json: binaries
	./bin/$(GOOS)/$(GOARCH)/$(PROJECT)$(GOEXE) expected-apis
	git diff --quiet apis.json || sh -c "git add apis.json ; git commit -m 'Update apis.json' apis.json"

.PHONY: update-modules
update-modules:
	zutano update-check --quiet --fail
	test -f go.mod || go mod init
	go mod edit \
		$(shell zutano go mod replacements)
	go get \
		$(shell zutano go mod latest \
			github.com/arangodb-managed/apis@remove-whitelist \
		)
	go mod tidy
