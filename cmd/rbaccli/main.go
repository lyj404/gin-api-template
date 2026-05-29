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
		fmt.Println("可用命令: create-admin, seed-resources, seed-menus, seed-dict")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create-admin":
		createAdmin()
	case "seed-resources":
		seedResources()
	case "seed-menus":
		seedMenus()
	case "seed-dict":
		seedDictData()
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
	fmt.Println()

	// 自动初始化基础数据
	bootstrap.BootDBOnly()
	seedResources()
	seedMenus()
	seedDictData()
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

		// 2.5 创建默认菜单（如果不存在）
		if err := createDefaultMenus(tx); err != nil {
			return fmt.Errorf("创建默认菜单失败: %w", err)
		}

		// 3. 创建/更新超级管理员角色
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
		}

		// 每次运行都重新绑定所有资源权限（确保新增资源同步到 super_admin）
		if err := syncAllResourcesToRole(tx, superAdminRole.ID); err != nil {
			return fmt.Errorf("绑定资源权限失败: %w", err)
		}

		// 每次运行都重新绑定所有菜单（确保新增菜单同步到 super_admin）
		if err := syncAllMenusToRole(tx, superAdminRole.ID); err != nil {
			return fmt.Errorf("绑定菜单失败: %w", err)
		}

		// 绑定全组织范围（仅首次创建时）
		orgScopesExist := tx.Model(&entity.RoleOrgScope{}).Where("role_id = ?", superAdminRole.ID).First(&entity.RoleOrgScope{}).Error == nil
		if !orgScopesExist {
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
		{Name: "role:bind-resource", Type: "api", Pattern: "/roles/:id/resources", Method: "POST", Description: "角色绑定资源"},
		{Name: "role:unbind-resource", Type: "api", Pattern: "/roles/:id/resources/:resourceId", Method: "DELETE", Description: "角色解绑资源"},
		{Name: "role:list-resources", Type: "api", Pattern: "/roles/:id/resources", Method: "GET", Description: "查看角色资源列表"},

		// API 资源 - 资源管理
		{Name: "resource:manage", Type: "api", Pattern: "/resources/*", Method: "*", Description: "资源管理"},

		// API 资源 - 组织管理
		{Name: "org:manage", Type: "api", Pattern: "/org-units/*", Method: "*", Description: "组织管理"},

		// API 资源 - 审计日志
		{Name: "audit:read", Type: "api", Pattern: "/audit-logs", Method: "GET", Description: "查看审计日志"},
		{Name: "audit:read:target", Type: "api", Pattern: "/audit-logs/target", Method: "GET", Description: "按目标查询审计日志"},
		{Name: "audit:read:time", Type: "api", Pattern: "/audit-logs/time", Method: "GET", Description: "按时间范围查询审计日志"},

		// API 资源 - 菜单管理
		{Name: "menu:read", Type: "api", Pattern: "/menus", Method: "GET", Description: "查看菜单列表"},
		{Name: "menu:read:detail", Type: "api", Pattern: "/menus/:id", Method: "GET", Description: "查看菜单详情"},
		{Name: "menu:read:tree", Type: "api", Pattern: "/menus/tree", Method: "GET", Description: "查看菜单树"},
		{Name: "menu:create", Type: "api", Pattern: "/menus", Method: "POST", Description: "创建菜单"},
		{Name: "menu:update", Type: "api", Pattern: "/menus/:id", Method: "PUT", Description: "更新菜单"},
		{Name: "menu:delete", Type: "api", Pattern: "/menus/:id", Method: "DELETE", Description: "删除菜单"},

		// API 资源 - 用户权限与菜单
		{Name: "user:permissions", Type: "api", Pattern: "/user/permissions", Method: "GET", Description: "获取用户权限"},
		{Name: "user:menus", Type: "api", Pattern: "/user/menus", Method: "GET", Description: "获取用户菜单"},

		// API 资源 - 仪表盘
		{Name: "dashboard:read", Type: "api", Pattern: "/dashboard/stats", Method: "GET", Description: "查看仪表盘统计"},
		{Name: "dashboard:audit-trend", Type: "api", Pattern: "/dashboard/audit-trend", Method: "GET", Description: "查看审计趋势"},

		// API 资源 - 字典管理
		{Name: "dict:read", Type: "api", Pattern: "/dict", Method: "GET", Description: "查看字典列表"},
		{Name: "dict:read:detail", Type: "api", Pattern: "/dict/:id", Method: "GET", Description: "查看字典详情"},
		{Name: "dict:create", Type: "api", Pattern: "/dict", Method: "POST", Description: "创建字典"},
		{Name: "dict:update", Type: "api", Pattern: "/dict/:id", Method: "PUT", Description: "更新字典"},
		{Name: "dict:delete", Type: "api", Pattern: "/dict/:id", Method: "DELETE", Description: "删除字典"},
		{Name: "dict:detail:read", Type: "api", Pattern: "/dict/:id/details", Method: "GET", Description: "查看字典详情列表"},
		{Name: "dict:detail:create", Type: "api", Pattern: "/dict/:id/details", Method: "POST", Description: "创建字典详情"},
		{Name: "dict:detail:update", Type: "api", Pattern: "/dict/:id/details/:detailId", Method: "PUT", Description: "更新字典详情"},
		{Name: "dict:detail:delete", Type: "api", Pattern: "/dict/:id/details/:detailId", Method: "DELETE", Description: "删除字典详情"},

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

func syncAllResourcesToRole(tx *gorm.DB, roleID uint64) error {
	// 清除旧绑定
	if err := tx.Where("role_id = ?", roleID).Delete(&entity.RoleResource{}).Error; err != nil {
		return err
	}

	// 重新绑定所有资源
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

func syncAllMenusToRole(tx *gorm.DB, roleID uint64) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&entity.RoleMenu{}).Error; err != nil {
		return err
	}

	var menus []entity.Menu
	if err := tx.Find(&menus).Error; err != nil {
		return err
	}

	for _, menu := range menus {
		rm := entity.RoleMenu{RoleID: roleID, MenuID: menu.ID}
		if err := tx.Create(&rm).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedResources() {
	fmt.Println("=== 资源数据初始化 ===")

	config.InitConfig()
	bootstrap.BootDBOnly()

	err := global.G_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建默认资源
		if err := createDefaultResources(tx); err != nil {
			return fmt.Errorf("创建默认资源失败: %w", err)
		}

		// 2. 查找 super_admin 角色，同步资源绑定
		var superAdminRole entity.Role
		if err := tx.Where("name = ?", "super_admin").First(&superAdminRole).Error; err == nil {
			if err := syncAllResourcesToRole(tx, superAdminRole.ID); err != nil {
				return fmt.Errorf("同步资源到 super_admin 角色失败: %w", err)
			}
			fmt.Println("已同步资源到 super_admin 角色")

			// 3. 确保 super_admin 角色有全组织范围（修复旧数据无 RoleOrgScope 的问题）
			var rootOrg entity.OrgUnit
			if err := tx.Where("name = ? AND parent_id IS NULL", "root").First(&rootOrg).Error; err == nil {
				var scopeCount int64
				tx.Model(&entity.RoleOrgScope{}).Where("role_id = ?", superAdminRole.ID).Count(&scopeCount)
				if scopeCount == 0 {
					if err := tx.Create(&entity.RoleOrgScope{
						RoleID:             superAdminRole.ID,
						OrgUnitID:          rootOrg.ID,
						IncludeDescendants: true,
					}).Error; err != nil {
						return fmt.Errorf("绑定组织范围到 super_admin 角色失败: %w", err)
					}
					fmt.Println("已绑定全组织范围到 super_admin 角色")
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("资源初始化失败: %v\n", err)
		os.Exit(1)
	}

	var count int64
	global.G_DB.Model(&entity.Resource{}).Count(&count)
	fmt.Printf("资源初始化成功，当前资源总数: %d\n", count)
}

func bindOrgScopeToRole(tx *gorm.DB, roleID, orgUnitID uint64, includeDescendants bool) error {
	roleOrgScope := entity.RoleOrgScope{
		RoleID:             roleID,
		OrgUnitID:          orgUnitID,
		IncludeDescendants: includeDescendants,
	}
	return tx.Create(&roleOrgScope).Error
}

func seedMenus() {
	fmt.Println("=== 菜单数据初始化 ===")

	config.InitConfig()
	bootstrap.BootDBOnly()

	err := global.G_DB.Transaction(func(tx *gorm.DB) error {
		if err := createDefaultMenus(tx); err != nil {
			return err
		}

		var superAdminRole entity.Role
		if err := tx.Where("name = ?", "super_admin").First(&superAdminRole).Error; err == nil {
			if err := syncAllMenusToRole(tx, superAdminRole.ID); err != nil {
				return fmt.Errorf("同步菜单到 super_admin 角色失败: %w", err)
			}
			fmt.Println("已同步菜单到 super_admin 角色")
		}

		return nil
	})
	if err != nil {
		fmt.Printf("菜单初始化失败: %v\n", err)
		os.Exit(1)
	}

	var count int64
	global.G_DB.Model(&entity.Menu{}).Count(&count)
	fmt.Printf("菜单初始化成功，当前菜单总数: %d\n", count)
}

func createDefaultMenus(tx *gorm.DB) error {
	type menuSeed struct {
		Name         string
		Path         string
		Icon         string
		OrderNum     int
		ResourceName string
	}

	seeds := []menuSeed{
		{Name: "仪表盘", Path: "/dashboard", Icon: "i-material-symbols:dashboard-outline", OrderNum: 10, ResourceName: "user:permissions"},
		{Name: "用户管理", Path: "/users", Icon: "i-material-symbols:group-outline", OrderNum: 20, ResourceName: "user:read"},
		{Name: "角色管理", Path: "/roles", Icon: "i-material-symbols:manage-accounts-outline", OrderNum: 30, ResourceName: "role:manage"},
		{Name: "菜单管理", Path: "/menus", Icon: "i-material-symbols:list-alt-outline", OrderNum: 40, ResourceName: "menu:read"},
		{Name: "组织管理", Path: "/orgs", Icon: "i-material-symbols:account-tree-outline", OrderNum: 50, ResourceName: "org:manage"},
		{Name: "资源管理", Path: "/resources", Icon: "i-material-symbols:shield-outline", OrderNum: 60, ResourceName: "resource:manage"},
		{Name: "字典管理", Path: "/dictionary", Icon: "i-material-symbols:book-outline", OrderNum: 65, ResourceName: "dict:read"},
		{Name: "审计日志", Path: "/audit-logs", Icon: "i-material-symbols:receipt-long-outline", OrderNum: 70, ResourceName: "audit:read"},
	}

	for _, s := range seeds {
		var existing entity.Menu
		if err := tx.Where("name = ? AND parent_id IS NULL", s.Name).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		var res entity.Resource
		if err := tx.Where("name = ?", s.ResourceName).First(&res).Error; err != nil {
			return fmt.Errorf("查找菜单关联资源 %s 失败: %w", s.ResourceName, err)
		}

		menu := entity.Menu{
			Name:      s.Name,
			ParentID:  nil,
			Path:      s.Path,
			Icon:      s.Icon,
			OrderNum:  s.OrderNum,
			IsVisible: true,
			Status:    "enabled",
		}
		if err := tx.Create(&menu).Error; err != nil {
			return err
		}

		mr := entity.MenuResource{MenuID: menu.ID, ResourceID: res.ID}
		if err := tx.Create(&mr).Error; err != nil {
			return err
		}
	}
	return nil
}

type dictSeed struct {
	Name  string
	Type  string
	Desc  string
	Items []dictItem
}

type dictItem struct {
	Label string
	Value string
	Sort  int
}

func seedDictData() {
	fmt.Println("=== 字典数据初始化 ===")

	config.InitConfig()
	bootstrap.BootDBOnly()

	seeds := []dictSeed{
		{
			Name: "菜单状态", Type: "menu_status", Desc: "菜单的启用/禁用状态",
			Items: []dictItem{
				{Label: "启用", Value: "enabled", Sort: 1},
				{Label: "禁用", Value: "disabled", Sort: 2},
			},
		},
		{
			Name: "资源类型", Type: "resource_type", Desc: "资源的类型分类",
			Items: []dictItem{
				{Label: "API — 接口权限", Value: "api", Sort: 1},
				{Label: "实体 — 数据权限", Value: "entity", Sort: 2},
			},
		},
		{
			Name: "审计日志操作", Type: "audit_action", Desc: "审计日志记录的操作类型",
			Items: []dictItem{
				{Label: "创建", Value: "create", Sort: 1},
				{Label: "更新", Value: "update", Sort: 2},
				{Label: "删除", Value: "delete", Sort: 3},
				{Label: "登录", Value: "login", Sort: 4},
				{Label: "登出", Value: "logout", Sort: 5},
				{Label: "注册", Value: "signup", Sort: 6},
				{Label: "绑定", Value: "bind", Sort: 7},
				{Label: "解绑", Value: "unbind", Sort: 8},
				{Label: "分配", Value: "assign", Sort: 9},
				{Label: "撤销", Value: "revoke", Sort: 10},
				{Label: "更新个人信息", Value: "update_profile", Sort: 11},
				{Label: "修改密码", Value: "change_password", Sort: 12},
			},
		},
		{
			Name: "审计日志目标类型", Type: "audit_target_type", Desc: "审计日志记录的目标对象类型",
			Items: []dictItem{
				{Label: "用户", Value: "user", Sort: 1},
				{Label: "角色", Value: "role", Sort: 2},
				{Label: "菜单", Value: "menu", Sort: 3},
				{Label: "资源", Value: "resource", Sort: 4},
				{Label: "组织", Value: "org_unit", Sort: 5},
				{Label: "角色资源", Value: "role_resource", Sort: 6},
				{Label: "角色组织范围", Value: "role_org_scope", Sort: 7},
				{Label: "用户角色", Value: "user_role", Sort: 8},
				{Label: "角色菜单", Value: "role_menu", Sort: 9},
				{Label: "菜单资源", Value: "menu_resource", Sort: 10},
			},
		},
	}

	for _, s := range seeds {
		var dict entity.SysDictionary
		err := global.G_DB.Where("type = ?", s.Type).First(&dict).Error
		if err == nil {
			// 更新已有字典值：删除旧值，重新插入
			global.G_DB.Where("dict_id = ?", dict.ID).Delete(&entity.SysDictionaryDetail{})
		} else {
			dict = entity.SysDictionary{
				Name:   s.Name,
				Type:   s.Type,
				Status: 1,
				Desc:   s.Desc,
			}
			if err := global.G_DB.Create(&dict).Error; err != nil {
				fmt.Printf("创建字典「%s」失败: %v\n", s.Name, err)
				continue
			}
		}

		for _, item := range s.Items {
			detail := entity.SysDictionaryDetail{
				DictID: dict.ID,
				Label:  item.Label,
				Value:  item.Value,
				Sort:   item.Sort,
				Status: 1,
			}
			if err := global.G_DB.Create(&detail).Error; err != nil {
				fmt.Printf("创建字典值「%s」失败: %v\n", item.Label, err)
			}
		}
		fmt.Printf("字典「%s」初始化完成 (%d 项)\n", s.Name, len(s.Items))
	}

	var count int64
	global.G_DB.Model(&entity.SysDictionary{}).Count(&count)
	fmt.Printf("字典数据初始化成功，当前字典总数: %d\n", count)
}
