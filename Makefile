# Project settings
BINARY := xizexcample
PKG := ./...
BUILD_DIR := bin
MAIN_PKG := .

GO ?= go

# 版本信息注入（可选）：
# 如需注入，先在 main 包中定义以下变量：
#   var version string
#   var commit  string
#   var date    string
# 然后将下方 LDFLAGS 替换为注入形式。
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")
# 默认仅压缩符号表，避免因未定义变量导致链接失败；
# 需要注入时可使用：
# LDFLAGS := -s -w -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'
LDFLAGS ?= -s -w

.DEFAULT_GOAL := build

.PHONY: build build-linux run test tidy clean proto docker-build docker-run fmt vet help

help:
	@echo "Common targets:"
	@echo "  make build         - 构建本地可执行文件 ($(BUILD_DIR)/$(BINARY))"
	@echo "  make build-linux   - 交叉编译 Linux/amd64 (CGO_DISABLED=0)"
	@echo "  make run           - 构建并运行"
	@echo "  make test          - 运行测试"
	@echo "  make tidy          - go mod tidy"
	@echo "  make clean         - 移除构建产物"
	@echo "  make proto         - 生成 Protobuf 代码 (调用 proto/generate_proto.sh)"
	@echo "  make docker-build  - 使用 Dockerfile 构建镜像"
	@echo "  make docker-run    - 运行镜像，映射 8080:8080"

build:
	mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)

# Cross-compile for linux/amd64
build-linux: export CGO_ENABLED=0
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: build

run: build
	./$(BUILD_DIR)/$(BINARY)

test:
	$(GO) test ./... -count=1

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt $(PKG)

vet:
	$(GO) vet $(PKG)

clean:
	rm -rf $(BUILD_DIR)

proto:
	cd proto && bash generate_proto.sh

docker-build:
	docker build -t $(BINARY):$(VERSION) .

docker-run:
	docker run --rm -p 8080:8080 $(BINARY):$(VERSION)
