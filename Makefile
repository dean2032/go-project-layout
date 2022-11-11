.PHONY: build clean

BIN="app"
GitCommit=$(shell git rev-parse --short HEAD || echo)
BuildTime=$(shell date +%Y-%m-%d.%H.%M.%S)
VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)

build:
	@echo "Building ${BIN}..."
	@go mod tidy
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GitCommit} -X main.BuildTime=${BuildTime}" -X main.BuildVersion=${VERSION} -o ${BIN} .
	@echo "built ${BIN} successed"

clean:
	@rm -r ${BIN} *.tgz *.tar.gz *.zip *.rpm dist/ log/
	@find . -name "*.log" | xargs rm

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.out -coverpkg=./...  
	go tool cover -html=coverage.out -o coverage.html
	
format:
	gocmt -d . -i -p
	golines -w -m 120 --no-reformat-tags .
