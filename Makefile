# 基础配置（简化+修复引号问题）
DYLD_FALLBACK_LIBRARY_PATH=/usr/lib
BINARY_NAME=main
DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty || echo "v1.0.0")
BUILD_TIME := $(shell date +%Y%m%d_%H%M%S)

# 启用CGO + 调整链接参数（兼容CGO静态链接）
LDFLAGS := -ldflags '-s -w -extldflags "-lpthread" -X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILD_TIME)"'
CGO_FLAG := CGO_ENABLED=1  # 核心：启用CGO

# 目标平台（格式：GOOS/GOARCH，清晰不易错）
TARGET_PLATFORMS := \
	linux/amd64/x86_64-linux-gnu-gcc \
	linux/arm64/aarch64-linux-gnu-gcc \
	darwin/amd64/clang \
	darwin/arm64/clang \
	windows/amd64/x86_64-w64-mingw32-gcc

proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/plugin/proto/*.proto
# ===================== 原有功能（完全保留） =====================
build-frontend:
	@echo "Building frontend..."
	cd ui && pnpm install && pnpm run build-only
	@echo "Copying frontend build to assets..."

build-backend:
	@echo "Building backend (local platform)..."
	CC=$(shell which gcc || which clang) $(CGO_FLAG) go build -o ${BINARY_NAME} ${LDFLAGS} main.go

build: build-frontend build-backend
	@echo "Build completed successfully!"

run: build
	./${BINARY_NAME}

clean:
	@echo "Cleaning all build artifacts..."
	go clean
	rm -f ${BINARY_NAME}
	@echo "Clean completed!"

dev:
	CC=$(shell which gcc || which clang) $(CGO_FLAG) go build -o ${BINARY_NAME} main.go && ./${BINARY_NAME}

frontend: build-frontend

# ===================== 修复后的多平台构建（核心） =====================
# 关键：所有Shell命令用单一行内循环，避免换行解析错误
.PHONY: build-all-platforms package
build-all-platforms: build-frontend
	@echo "=== Start building all platforms (CGO enabled) ==="
	@for platform in $(TARGET_PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d '/' -f1); \
		GOARCH=$$(echo $$platform | cut -d '/' -f2); \
		CC=$$(echo $$platform | cut -d '/' -f3); \
		\
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; else EXT=""; fi; \
		OUTPUT_DIR=$(DIST_DIR)/$$GOOS"_"$$GOARCH; \
		OUTPUT_FILE=$$OUTPUT_DIR/$(BINARY_NAME)$$EXT; \
		mkdir -p $$OUTPUT_DIR; \
		echo "--- Building $$GOOS/$$GOARCH (CC=$$CC) ---"; \
		\
		CC=$$CC $(CGO_FLAG) GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $$OUTPUT_FILE main.go; \
		echo "✅ $$GOOS/$$GOARCH built to $$OUTPUT_FILE"; \
	done
	@echo "=== All platforms build completed! ==="

# 打包产物（同样用行内循环，避免语法错误）
package: build-all-platforms
	@echo "=== Packaging all artifacts ==="
	@for platform in $(TARGET_PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d '/' -f1); \
		GOARCH=$$(echo $$platform | cut -d '/' -f2); \
		PLATFORM_DIR=$$GOOS"_"$$GOARCH; \
		ARCHIVE_NAME=$(BINARY_NAME)-$(VERSION)-$$PLATFORM_DIR; \
		mkdir -p $(DIST_DIR); \
		cd $(DIST_DIR); \
		if [ "$$GOOS" = "windows" ]; then \
			zip -q -r $$ARCHIVE_NAME.zip $$PLATFORM_DIR/* 2>/dev/null; \
		else \
			tar -zcf $$ARCHIVE_NAME.tar.gz $$PLATFORM_DIR/* 2>/dev/null; \
		fi; \
		echo "✅ Packaged $(DIST_DIR)/$$ARCHIVE_NAME.$$( [ "$$GOOS" = "windows" ] && echo "zip" || echo "tar.gz" )"; \
	done
	@echo "=== Packaging completed! ==="