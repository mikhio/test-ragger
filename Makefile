BIN := bin/test-ragger
PKG := ./cmd/test-ragger

# Use public Go proxy to avoid corporate proxy issues
GOPROXY := https://proxy.golang.org,direct

.PHONY: all deps build run clean ingest search fmt vet tidy test docker-up docker-down help

all: build

# -------- Setup and dependencies --------
deps:
	GOPROXY=$(GOPROXY) go mod download
	GOPROXY=$(GOPROXY) go mod tidy

setup: deps
	@mkdir -p html bin
	@echo "✅ Project setup complete"
	@echo "📝 Don't forget to set OPENAI_API_KEY in your environment"

# -------- Build --------
build:
	@mkdir -p bin
	GOPROXY=$(GOPROXY) go build -o $(BIN) $(PKG)

build-clean: clean build

# -------- Run commands --------
run: build
	./$(BIN)

# Default parameters
DIR ?= ./html
MODEL ?= text-embedding-3-small
QDRANT ?= localhost:6334
K ?= 5
Q ?=
LANG ?=

ingest: build
	@echo "🔄 Running ingest mode..."
	./$(BIN) -mode=ingest -dir=$(DIR) -qdrant=$(QDRANT) -model=$(MODEL)

search: build
	@[ -n "$(Q)" ] || (echo "❌ Q is required (query). Usage: make search Q='your query'" && exit 1)
	@echo "🔍 Searching for: $(Q)"
	./$(BIN) -mode=search -q="$(Q)" -k=$(K) -qdrant=$(QDRANT) $(if $(LANG),-lang=$(LANG),)

# -------- Docker commands --------
docker-up:
	@echo "🐳 Starting Qdrant..."
	docker-compose up -d
	@echo "✅ Qdrant is running at http://localhost:6333/dashboard"

docker-down:
	@echo "🛑 Stopping Qdrant..."
	docker-compose down

docker-logs:
	docker-compose logs -f qdrant

docker-clean:
	docker-compose down -v
	@echo "🧹 Qdrant data cleared"

# -------- Development --------
test:
	GOPROXY=$(GOPROXY) go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	GOPROXY=$(GOPROXY) go mod tidy

lint: fmt vet
	@echo "✅ Code linting complete"

# -------- Cleanup --------
clean:
	rm -rf bin

clean-all: clean docker-clean
	go clean -modcache
	@echo "🧹 Full cleanup complete"

# -------- Help --------
help:
	@echo "📖 Available commands:"
	@echo ""
	@echo "🔧 Setup:"
	@echo "  make setup        - Initial project setup"
	@echo "  make deps         - Download and tidy dependencies"
	@echo ""
	@echo "🏗️  Build:"
	@echo "  make build        - Build the application"
	@echo "  make build-clean  - Clean build"
	@echo ""
	@echo "🚀 Run:"
	@echo "  make ingest [DIR=./html] [MODEL=text-embedding-3-small]"
	@echo "  make search Q='query' [K=5] [LANG=ru]"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-up    - Start Qdrant"
	@echo "  make docker-down  - Stop Qdrant"
	@echo "  make docker-logs  - Show Qdrant logs"
	@echo "  make docker-clean - Stop and clear Qdrant data"
	@echo ""
	@echo "🧹 Development:"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Format and vet code"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make clean-all    - Full cleanup"
	@echo ""
	@echo "💡 Examples:"
	@echo "  OPENAI_API_KEY=sk-... make ingest"
	@echo "  OPENAI_API_KEY=sk-... make search Q='машинное обучение' K=10"


