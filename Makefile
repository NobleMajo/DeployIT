DISPLAY_NAME := DeployIT
SHORT_NAME := dit
MAIN_GO := main.go

VERSION := $(shell .github/scripts/latest-dev-version.sh)
COMMIT := $(shell git rev-parse --short HEAD)
PORT ?= 8080
-include .env
export

.DEFAULT_GOAL = help
.MAIN = help
.PHONY: help # prints the help message
help:
	@echo "$(DISPLAY_NAME) version $(VERSION), build $(COMMIT)"
	@echo ""
	@echo "Available make targets:"
	@grep -hE '^\.PHONY: [a-zA-Z0-9_.-]+ #' $(MAKEFILE_LIST) \
		| sed 's/^\.PHONY: //' \
		| awk -F' #' '{ gsub(/^[ \t]+|[ \t]+$$/, "", $$2); \
			printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }'

.PHONY: run # uses go to run the main program (MAIN_GO)
run:
	@go run $(MAIN_GO)

.PHONY: build # uses go to build the app with build args
build:
	@touch .env
	go build \
		-ldflags="$(shell $(MAKE) -s buildflags)" \
		-o bin \
		$(MAIN_GO)
	chmod +x bin

.PHONY: buildflags # builds the build flags for the go build command
buildflags:
	@echo "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.DisplayName=$(DISPLAY_NAME) -X main.ShortName=$(SHORT_NAME)"

.PHONY: short_name # prints the short name
short_name:
	@echo $(SHORT_NAME)

.PHONY: clean # cleans up the tmp, build and docker cache
clean:
	@echo "Clearing temporary files..."
	rm -rf ./tmp ./bin
	@if command -v go 2>&1 >/dev/null; then \
		echo "cleanup go..."; \
		go clean -cache -fuzzcache; \
	fi
	@if command -v docker 2>&1 >/dev/null; then \
		echo "cleanup docker..."; \
		docker compose down --remove-orphans --rmi all; \
		docker image prune -f; \
	fi
	@echo "cleanup done!"
	@echo "WARNING: the .env file still exists!"

.PHONY: update # updates go dependencies
update:
	go get -t -u ./...

.PHONY: test # runs all tests without coverage
test:
	go vet ./...
	go test -failfast ./...
	make -s build
	./bin

.PHONY: dev # starts the go bin in watch mode
dev:
	@go install github.com/air-verse/air@v1
	@air

.PHONY: docker # starts an interactive bash inside a dev docker container
docker:
	@touch .env
	@mkdir -p ./tmp
	@docker compose run \
	    --build --rm -P -it \
	    --name $(SHORT_NAME)-local-bash \
		local

.PHONY: tag-patch # tag a patch release
tag-patch:
	@.github/scripts/release-git-tag.sh patch

.PHONY: tag-minor # tag a minor release
tag-minor:
	@.github/scripts/release-git-tag.sh minor

.PHONY: tag-major # tag a major release
tag-major:
	@.github/scripts/release-git-tag.sh major

