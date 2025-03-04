package pagination

import (
	"math"
	"strings"

	"gorm.io/gorm"
)

// PaginationBuilder 分页查询构造器
// 使用构建器模式，支持链式调用来构建复杂的分页查询
type PaginationBuilder struct {
	db          *gorm.DB    // 数据库连接实例
	page        int         // 当前页
	pageSize    int         // 每页数量
	model       interface{} // 模型
	preloads    []string    // 添加预加载
	joins       []string    // 添加连接查询
	selects     []string    // 添加查询字段
	conditions  []Condition // 添加查询条件
	orderFields []string    // 排序字段
	groups      []string    // 添加分组字段
	havings     []Condition // 添加HAVING条件
}

// NewPaginationBuilder 创建分页构造器的新实例
// 参数:
//   - db: *gorm.DB GORM数据库连接实例
//
// 返回:
//   - *PaginationBuilder 新的分页构造器实例
func NewPaginationBuilder(db *gorm.DB) *PaginationBuilder {
	return &PaginationBuilder{
		db:       db,
		page:     1,
		pageSize: 10,
	}
}

// Build 执行分页查询并构建结果
// 参数:
//   - result: interface{} 查询结果的接收对象指针
//
// 返回:
//   - *Pagination 分页查询结果
//   - error 查询过程中的错误
//
// 示例:
//
//	var users []User
//	pagination, err := builder.Build(&users)
func (p *PaginationBuilder) Build(result interface{}) (*Pagination, error) {
	var total int64
	query := p.db.Model(p.model)

	// 查询条件
	for _, condition := range p.conditions {
		query = query.Where(condition.Query, condition.Args...)
	}

	// 连接查询
	for _, join := range p.joins {
		query = query.Joins(join)
	}

	// 查询字段
	if len(p.selects) > 0 {
		query = query.Select(strings.Join(p.selects, ","))
	}

	// 添加分组
	for _, group := range p.groups {
		query = query.Group(group)
	}

	// 添加HAVING条件
	for _, having := range p.havings {
		query = query.Having(having.Query, having.Args...)
	}

	// 克隆查询以计算总记录数
	countQuery := query
	if len(p.groups) > 0 {
		// 如果有分组，则计算分组后的记录数
		countQuery = countQuery.Session(&gorm.Session{})
		if err := countQuery.Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		// 无分组时的普通计数
		if err := countQuery.Count(&total).Error; err != nil {
			return nil, err
		}
	}

	// 添加排序
	for _, field := range p.orderFields {
		query = query.Order(field)
	}

	// 分页查询
	offset := (p.page - 1) * p.pageSize
	query = query.Offset(offset).Limit(p.pageSize)

	// 预加载
	for _, preload := range p.preloads {
		query = query.Preload(preload)
	}

	// 执行查询
	if err := query.Find(result).Error; err != nil {
		return nil, err
	}

	// 构建分页响应
	totalPage := int(math.Ceil(float64(total) / float64(p.pageSize)))
	return &Pagination{
		Page:      p.page,
		PageSize:  p.pageSize,
		Total:     total,
		TotalPage: totalPage,
		Data:      result,
	}, nil
}
