export CGO_CPPFLAGS=${CPPFLAGS}
export CGO_CFLAGS=${CFLAGS}
export CGO_CXXFLAGS=${CXXFLAGS}
export CGO_LDFLAGS=${LDFLAGS}
GOFLAGS ?= -buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw

# update go dependencies
update:
	go get ./cmd
	go mod tidy
	go mod vendor

mock:
	@mockery --all

# run linter
lint:
	golangci-lint run ./...

# run linter and fix issues if possible
lintfix:
	golangci-lint run --fix ./...

# run unit tests
test:
	@go test -coverprofile=cover.out ./...
	@go tool cover -func=cover.out
	-@rm -f cover.out

# run emm, note: make doesn't understand exit code 130 and sets it == 1
run:
	@go run ./cmd || exit 0

install: build
	@mv ./emm ${HOME}/go/bin/
	@echo "emm has been installed."

# build emm
build:
	go build -v -o emm ./cmd

