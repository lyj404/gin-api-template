# 定义伪目标
.PHONY: all build run clean swagger create-admin seed-resources seed-menus seed clean-logs

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

# 初始化系统资源到数据库
seed-resources:
	$(GO) run ./cmd/rbaccli/main.go seed-resources

# 创建系统管理员默认菜单
seed-menus:
	$(GO) run ./cmd/rbaccli/main.go seed-menus

# 初始化所有基础数据（资源 + 菜单 + 字典）
seed: seed-resources seed-menus seed-dict

# 初始化系统字典数据（菜单状态、资源类型等）
seed-dict:
	$(GO) run ./cmd/rbaccli/main.go seed-dict

