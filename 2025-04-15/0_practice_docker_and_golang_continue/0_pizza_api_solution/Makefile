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
	go mod tidy
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/app" ./cmd/app

build_debug:
	go mod tidy
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS_DEBUG}" -o "${BUILD_DIR}/app" ./cmd/app

run: build
	"${BUILD_DIR}/app"

docker_alpine:
	docker build --file Dockerfile_alpine -t "${IMAGE_NAME_ALPINE}:${VERSION}" .
	docker run --env-file .env -p 8080:8080 "${IMAGE_NAME_ALPINE}:${VERSION}"

docker_distroless:
	docker build --file Dockerfile_distroless -t "${IMAGE_NAME_DISTROLESS}:${VERSION}" .
	docker run --env-file .env -p 8080:8080 "${IMAGE_NAME_DISTROLESS}:${VERSION}"

docker_scratch:
	docker build --file Dockerfile_scratch -t "${IMAGE_NAME_SCRATCH}:${VERSION}" .
	docker run --env-file .env -p 8080:8080 "${IMAGE_NAME_SCRATCH}:${VERSION}"

docker_bookworm:
	docker build --file Dockerfile_bookworm -t "${IMAGE_NAME_BOOKWORM}:${VERSION}" .
	docker run --env-file .env -p 8080:8080 "${IMAGE_NAME_BOOKWORM}:${VERSION}"

docker_alpine_debug:
	docker build --file Dockerfile_alpine_debug -t "${IMAGE_NAME_ALPINE}:${VERSION}" .
	docker run --env-file .env -p 8080:8080 -p 2345:2345 "${IMAGE_NAME_ALPINE}:${VERSION}"
