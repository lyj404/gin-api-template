package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/lyj404/gin-api-template/bootstrap"
	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain/entity"
	"github.com/lyj404/gin-api-template/global"
	"github.com/lyj404/gin-api-template/util"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("可用命令: create-admin")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create-admin":
		createAdmin()
	default:
		fmt.Printf("未知命令: %s\n", command)
		os.Exit(1)
	}
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(password)
}

func createAdmin() {
	fmt.Println("=== 系统管理员初始化 ===")

	config.InitConfig()
	bootstrap.BootDBOnly()

	reader := bufio.NewReader(os.Stdin)

	var email, password string

	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" {
		email = adminEmail
		fmt.Printf("使用环境变量 ADMIN_EMAIL: %s\n", email)
	} else {
		fmt.Print("输入管理员邮箱: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	if email == "" {
		fmt.Println("邮箱不能为空，请设置 ADMIN_EMAIL 环境变量或手动输入")
		os.Exit(1)
	}

	if adminPassword := os.Getenv("ADMIN_PASSWORD"); adminPassword != "" {
		password = adminPassword
		fmt.Println("使用环境变量 ADMIN_PASSWORD")
	} else {
		password = readPassword("输入管理员密码: ")
	}

	if password == "" {
		fmt.Println("密码不能为空，请设置 ADMIN_PASSWORD 环境变量或手动输入")
		os.Exit(1)
	}

	if os.Getenv("ADMIN_PASSWORD") == "" {
		confirmPassword := readPassword("确认管理员密码: ")

		if password != confirmPassword {
			fmt.Println("两次密码不一致")
			os.Exit(1)
		}
	}

	err := createSystemAdmin(email, password)
	if err != nil {
		fmt.Printf("创建管理员失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("系统管理员创建成功！")
}

func createSystemAdmin(email, password string) error {
	return global.G_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建根组织节点（如果不存在）
		var rootOrg entity.OrgUnit
		result := tx.Where("name = ? AND parent_id IS NULL", "root").First(&rootOrg)
		if result.Error == gorm.ErrRecordNotFound {
			rootOrg = entity.OrgUnit{
				Name:     "root",
				ParentID: nil,
				Path:     "/0",
				Level:    0,
			}
			if err := tx.Create(&rootOrg).Error; err != nil {
				return fmt.Errorf("创建根组织失败: %w", err)
			}
		}

		// 2. 创建默认资源（如果不存在）
		if err := createDefaultResources(tx); err != nil {
			return fmt.Errorf("创建默认资源失败: %w", err)
		}

		// 3. 创建超级管理员角色（如果不存在）
		var superAdminRole entity.Role
		result = tx.Where("name = ?", "super_admin").First(&superAdminRole)
		if result.Error == gorm.ErrRecordNotFound {
			superAdminRole = entity.Role{
				Name:        "super_admin",
				Description: "超级管理员",
				IsSystem:    true,
			}
			if err := tx.Create(&superAdminRole).Error; err != nil {
				return fmt.Errorf("创建超级管理员角色失败: %w", err)
			}

			// 绑定所有资源权限
			if err := bindAllResourcesToRole(tx, superAdminRole.ID); err != nil {
				return fmt.Errorf("绑定资源权限失败: %w", err)
			}

			// 绑定全组织范围
			if err := bindOrgScopeToRole(tx, superAdminRole.ID, rootOrg.ID, true); err != nil {
				return fmt.Errorf("绑定组织范围失败: %w", err)
			}
		}

		// 4. 创建管理员用户
		encryptedPassword, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("密码加密失败: %w", err)
		}

		adminUser := entity.User{
			Name:     "系统管理员",
			Email:    email,
			PassWord: string(encryptedPassword),
		}

		if err := tx.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}

		// 5. 绑定超级管理员角色
		userRole := entity.UserRole{
			UserID:    adminUser.ID,
			RoleID:    superAdminRole.ID,
			OrgUnitID: rootOrg.ID,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return fmt.Errorf("绑定角色失败: %w", err)
		}

		// 6. 记录审计日志
		auditLog := entity.AuditLog{
			OperatorID:   adminUser.ID,
			OperatorName: adminUser.Name,
			Action:       "create",
			TargetType:   "user",
			TargetID:     adminUser.ID,
			Description:  "创建系统管理员并绑定超级管理员角色",
		}

		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("记录审计日志失败: %w", err)
		}

		return nil
	})
}

func createDefaultResources(tx *gorm.DB) error {
	defaultResources := []entity.Resource{
		// API 资源 - 用户管理
		{Name: "user:create", Type: "api", Pattern: "/users", Method: "POST", Description: "创建用户"},
		{Name: "user:read", Type: "api", Pattern: "/users", Method: "GET", Description: "查看用户列表"},
		{Name: "user:read:detail", Type: "api", Pattern: "/users/:id", Method: "GET", Description: "查看用户详情"},
		{Name: "user:update", Type: "api", Pattern: "/users/:id", Method: "PUT", Description: "更新用户"},
		{Name: "user:delete", Type: "api", Pattern: "/users/:id", Method: "DELETE", Description: "删除用户"},

		// API 资源 - 角色管理
		{Name: "role:manage", Type: "api", Pattern: "/roles/*", Method: "*", Description: "角色管理"},

		// API 资源 - 资源管理
		{Name: "resource:manage", Type: "api", Pattern: "/resources/*", Method: "*", Description: "资源管理"},

		// API 资源 - 组织管理
		{Name: "org:manage", Type: "api", Pattern: "/org-units/*", Method: "*", Description: "组织管理"},

		// API 资源 - 审计日志
		{Name: "audit:read", Type: "api", Pattern: "/audit-logs", Method: "GET", Description: "查看审计日志"},

		// API 资源 - 菜单管理
		{Name: "menu:read", Type: "api", Pattern: "/menus", Method: "GET", Description: "查看菜单列表"},
		{Name: "menu:read:tree", Type: "api", Pattern: "/menus/tree", Method: "GET", Description: "查看菜单树"},
		{Name: "menu:create", Type: "api", Pattern: "/menus", Method: "POST", Description: "创建菜单"},
		{Name: "menu:update", Type: "api", Pattern: "/menus/:id", Method: "PUT", Description: "更新菜单"},
		{Name: "menu:delete", Type: "api", Pattern: "/menus/:id", Method: "DELETE", Description: "删除菜单"},

		// API 资源 - 用户权限与菜单
		{Name: "user:permissions", Type: "api", Pattern: "/user/permissions", Method: "GET", Description: "获取用户权限"},
		{Name: "user:menus", Type: "api", Pattern: "/user/menus", Method: "GET", Description: "获取用户菜单"},

		// 实体资源
		{Name: "entity:all", Type: "entity", Pattern: "*", Entity: "*", Action: "*", Description: "所有实体权限"},
	}

	for _, resource := range defaultResources {
		var existing entity.Resource
		result := tx.Where("name = ?", resource.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := tx.Create(&resource).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func bindAllResourcesToRole(tx *gorm.DB, roleID uint) error {
	var resources []entity.Resource
	if err := tx.Find(&resources).Error; err != nil {
		return err
	}

	for _, resource := range resources {
		roleResource := entity.RoleResource{
			RoleID:     roleID,
			ResourceID: resource.ID,
			IsRead:     true,
			IsWrite:    true,
		}
		if err := tx.Create(&roleResource).Error; err != nil {
			return err
		}
	}
	return nil
}

func bindOrgScopeToRole(tx *gorm.DB, roleID, orgUnitID uint, includeDescendants bool) error {
	roleOrgScope := entity.RoleOrgScope{
		RoleID:             roleID,
		OrgUnitID:          orgUnitID,
		IncludeDescendants: includeDescendants,
	}
	return tx.Create(&roleOrgScope).Error
}
