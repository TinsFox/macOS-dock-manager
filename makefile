# 变量定义
BINARY_NAME=dock-manager
GO=go
GOCLEAN=$(GO) clean
GOBUILD=$(GO) build
GOMOD=$(GO) mod
GOGET=$(GO) get
BUILD_DIR=build

# 默认目标
.DEFAULT_GOAL := build-and-run

# 检查并初始化 go.mod
init-mod:
	@if [ ! -f go.mod ]; then \
		echo "Initializing project..."; \
		$(GOMOD) init dock-manager; \
	fi
	@echo "Installing dependencies..."
	$(GOGET) github.com/AlecAivazis/survey/v2

# 创建构建目录
create-build-dir:
	@mkdir -p $(BUILD_DIR)

# 构建项目
build: init-mod create-build-dir
	@echo "Building..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)

# 清理项目
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# 运行程序
run:
	@echo "Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 构建并运行（主要使用目标）
build-and-run: build run

# 全部执行（构建、运行）
all: build-and-run

# 帮助信息
help:
	@echo "Available commands:"
	@echo "  make          - Build and run the application"
	@echo "  make build    - Only build the application"
	@echo "  make clean    - Clean build files"
	@echo "  make run      - Run the application (requires previous build)"
	@echo "  make all      - Build and run"
	@echo "  make help     - Show this help message"

.PHONY: init-mod build clean run build-and-run all help
