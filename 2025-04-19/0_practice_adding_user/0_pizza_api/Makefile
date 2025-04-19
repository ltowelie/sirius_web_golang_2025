-include .env
export

PROJECT_PKG := pizza_api
BUILD_DIR := "build"

VERSION ?=$(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
# inject build info
LDFLAGS += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}
LDFLAGS_DEBUG += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}

.PHONY: build swagger migrations install_tools docker_alpine

swagger:
	swag init -g cmd/app/main.go -o api_docs -parseDependency

install_tools:
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	
build:
	rm -rf "${BUILD_DIR}/app"
	go mod tidy
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/app" ./cmd/app

build_debug:
	rm -rf "${BUILD_DIR}/app"
	go mod tidy
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS_DEBUG}" -o "${BUILD_DIR}/app" ./cmd/app

run: build
	"${BUILD_DIR}/app"

docker_compose_production:
	docker compose --profile production up

docker_compose_develop_air:
	rm -rf "${BUILD_DIR}/app"
	docker compose --profile develop_air up

docker_compose_develop_air_debug:
	rm -rf "${BUILD_DIR}/app"
	docker compose --profile develop_air_debug up
