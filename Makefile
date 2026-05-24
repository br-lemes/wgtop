.PHONY: build clean linux-amd64 linux-arm64 release test version

TARGET := $(notdir $(shell go list -m 2>/dev/null))
ifeq ($(TARGET),)
	TARGET := $(notdir $(CURDIR))
endif

build: test
	@go build -ldflags "-s -w"

clean:
	$(RM) $(TARGET) $(TARGET)-linux-amd64 $(TARGET)-linux-arm64

linux-amd64: test
	@CGO_ENABLED=0 \
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(TARGET)-$@

linux-arm64: test
	@CGO_ENABLED=0 \
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(TARGET)-$@

release: version linux-amd64 linux-arm64
	@go run ./tools/release/main.go

test:
	@go test ./...

version: test
	@go run ./tools/version/main.go
