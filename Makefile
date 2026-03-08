VERSION := dev-$(shell date +%Y-%m-%d_%H:%M:%S)
INSTALL_DIR := ~/pink-tools/pink-elevenlabs

build:
	go build -ldflags="-X main.version=$(VERSION)" -o pink-elevenlabs .

install: build
	cp pink-elevenlabs $(INSTALL_DIR)/pink-elevenlabs

setup:
	git config core.hooksPath .githooks
