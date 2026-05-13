# 定义伪目标
.PHONY: all build run clean swagger create-admin seed-menus clean-logs

# 项目名称
PROJECT_NAME := gin-api-template

# Go 命令
GO := go

# 操作系统检测
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    BINARY_NAME := $(PROJECT_NAME).exe
    RM := if exist bin rmdir /s /q bin
    RM_LOGS := if exist logs rmdir /s /q logs
else
    BINARY_NAME := $(PROJECT_NAME)
    RM := rm -rf bin
    RM_LOGS := rm -rf ./logs/*
endif

# 默认目标
all: run

# 构建项目
build:
	$(GO) build -ldflags "-s -w" -o bin/$(BINARY_NAME) ./cmd

# 运行项目
run:
	cd cmd && $(GO) run .

# 清理构建文件
clean:
	@$(RM)

# 清除日志文件
clean-logs:
	@$(RM_LOGS)

swagger:
	swag init -g ./cmd/main.go --parseDependency

# 创建系统管理员
create-admin:
	$(GO) run ./cmd/rbaccli/main.go create-admin

# 创建系统管理员默认菜单
seed-menus:
	$(GO) run ./cmd/rbaccli/main.go seed-menus

