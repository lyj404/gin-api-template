# 定义伪目标
.PHONY: all build run clean swagger

# 项目名称
PROJECT_NAME := gin-api-template

# Go 命令
GO := go

# 动态检测操作系统和架构
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# 构建目标目录
BUILD_DIR := ./bin
BINARY_NAME := $(PROJECT_NAME)

# 默认目标
all: run

# 构建项目
build:
	@echo "🚀 构建项目中..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(PROJECT_NAME)
	@echo "✅ 构建完成！二进制文件位于 $(BUILD_DIR)/$(BINARY_NAME)"

# 运行项目
run:
	@echo "🏃 运行项目中..."
	$(GO) run ./cmd/main.go

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -rf $(BUILD_DIR)
	@echo "✅ 清理完成！"

# 清除日志文件
clean-logs:
	@echo "🧹 清理日志文件..."
	@rm -rf ./logs/*
	@echo "✅ 成功清理日志文件！"

swagger:
	@echo "📚 生成 Swagger 文档..."
	swag init -g ./cmd/main.go --parseDependency
	@echo "✅ Swagger 文档生成完成！"