BIN := bin/test-ragger
PKG := ./cmd/test-ragger

.PHONY: all deps build run clean ingest search fmt vet tidy

all: build

deps:
	GOPROXY=https://proxy.golang.org,direct go mod tidy

build:
	@mkdir -p bin
	go build -o $(BIN) $(PKG)

run: build
	./$(BIN)

clean:
	rm -rf bin

# -------- Commands with params --------
# Usage:
#   OPENAI_API_KEY=... make ingest DIR=./html MODEL=text-embedding-3-small QDRANT=localhost:6334
#   OPENAI_API_KEY=... make search Q="your query" K=5 QDRANT=localhost:6334

DIR ?= ./html
MODEL ?= text-embedding-3-small
QDRANT ?= localhost:6334
K ?= 5
Q ?=

ingest: build
	./$(BIN) -mode ingest -dir $(DIR) -qdrant $(QDRANT) -model $(MODEL)

search: build
	@[ -n "$(Q)" ] || (echo "Q is required (query)" && exit 1)
	./$(BIN) -mode search -q "$(Q)" -k $(K) -qdrant $(QDRANT)

# -------- Hygiene --------
fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	GOPROXY=https://proxy.golang.org,direct go mod tidy


