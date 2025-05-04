

GOOS := $(shell go env GOOS)
NAME := find

ifeq ($(GOOS), linux)
	BINARY_NAME := $(NAME)
endif

ifeq ($(GOOS), windows)
	BINARY_NAME := $(NAME).exe
endif

all: build

build:
	mkdir -p bin
	go build -o bin/$(NAME).exe src/main.go

run:
	go run .

clean:
	rm -rf bin/

.PHONY: all build run clean
