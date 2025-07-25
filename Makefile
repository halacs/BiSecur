DIST=dist/
APPNAME=halsecur

GOLANGCILINT_VERSION=v2.1.5
GOSEC_VERSION=v2.22.3
VULNCHECK_VERSION=latest

all: clean build

env:
	mkdir -p ${DIST}

clean:
	rm -rf ${DIST}

lint-env:
	( which gosec &>/dev/zero && gosec --version | grep -qs $(GOSEC_VERSION) ) || go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	( which golangci-lint &>/dev/zero && golangci-lint --version | grep -qs $(GOLANGCILINT_VERSION) ) || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCILINT_VERSION)
	( which govulncheck &>/dev/zero ) || go install golang.org/x/vuln/cmd/govulncheck@$(VULNCHECK_VERSION)

lint: lint-env
	golangci-lint --timeout 10m -v run ./...
	gosec ./...
	govulncheck ./...

lint-fix: lint-env
	golangci-lint run -v --fix ./...

test: test-short
	go test ${VENDOR} ./...

test-short:
	go test ${VENDOR} -race -short

build: env
	CGO_ENABLED=0 go build -ldflags "-X 'bisecur/version.Version=?version?' -X 'bisecur/version.BuildDate=?date?'" -v -o ${DIST}${APPNAME} .

build-docker: env
	docker build --build-arg VERSION=$(shell git describe --tags --always) -t bisecur/halsecur:latest .
