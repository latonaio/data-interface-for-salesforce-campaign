.DEFAULT_GOAL := help

GO_SRCS := $(shell find . -type f -name '*.go')

.PHONY: count-go
count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

.PHONY: docker-build
docker-build: $(GO_SRCS) ## docker image を build します。
	bash ./scripts/build.sh

.PHONY: go-build
go-build: $(GO_SRCS) ## go build を実行します。
	go build ./

.PHONY: go-test
go-test: $(GO_SRCS) ## go test を実行します。
	go test -v

# See "Self-Documented Makefile" article
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
