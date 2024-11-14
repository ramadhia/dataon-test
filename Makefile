APP_NAME=backend
VERSION_VAR=main.Version
VERSION=$(shell git describe --tags)
MOCKERY_VERSION=2.16.0
OS=Linux
ARCH=x86_64

build: dep
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X ${VERSION_VAR}=${VERSION}" -a -installsuffix nocgo -o ./bin ./...

dep:
	@echo ">> Downloading Dependencies"
	@go mod download

dep-lint:
	@echo ">> Downloading golangci-lint"
	@curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b $$(go env GOPATH)/bin v1.50.1

lint:
	@echo ">> Running Linter"
	@$$(go env GOPATH)/bin/golangci-lint run

lint-new:
	@echo ">> Running Linter on new/updated files"
	@$$(go env GOPATH)/bin/golangci-lint run --new-from-rev=HEAD~1

docker:
	@echo ">> Building Docker Image"
	@docker build -t ${APP_NAME}:latest .

hooks:
	@echo ">> Installing git hooks"
	git config core.hooksPath .githooks

run-api: dep
	@echo ">> Running API Server"
	@env $$(cat .env | xargs) go run github.com/ramadhia/mnc-test/cmd/cli server

run-worker: dep
	@echo ">> Running Worker Server"
	@env $$(cat .env | xargs) go run github.com/ramadhia/mnc-test/cmd/cli worker

run-rabbit:
	@echo ">> Starting RabbitMQ"
	#docker run -d --hostname my-rabbit --name some-rabbit -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:latest
	docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

migrate: dep
	@echo ">> Running DB migration"
	@env $$(cat .env | xargs) go run github.com/ramadhia/mnc-test/cmd/cli migrate

dep-mock:
ifeq (, $(shell which mockery))
	@echo ">> Download mockery"
	@curl -sSfL "https://github.com/vektra/mockery/releases/download/v${MOCKERY_VERSION}/mockery_${MOCKERY_VERSION}_${OS}_${ARCH}.tar.gz" | tar -zx -C $$(go env GOPATH)/bin mockery
endif
