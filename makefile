.PHONY: build_linux build_macos build_windows build_all

GOROOT=$(shell go env GOROOT)
GO=$(GOROOT)/bin/go
GOBUILD=GOROOT=$(GOROOT) CGO_ENABLED=0 $(GO) build -o

build_linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) ./build/ceobebot-qqchanel_linux-amd64 .

build_macos:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) ./build/ceobebot-qqchanel_darwin-arm64 .

build_macos_amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) ./build/ceobebot-qqchanel_darwin-amd64 .

build_windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) ./build/ceobebot-qqchanel_windows-amd64.exe .

build_all: build_linux build_macos build_macos_amd64 build_windows
