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
	preloads    []string    // 预加载关系
	joins       []string    // 连接查询
	selects     []string    // 查询字段
	distinct    bool        // 是否启用 DISTINCT
	conditions  []Condition // 查询条件
	orderFields []string    // 排序字段
	groups      []string    // 分组字段
	havings     []Condition // HAVING条件
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

// Joins 添加JOIN查询
func (p *PaginationBuilder) Joins(query string) *PaginationBuilder {
	p.joins = append(p.joins, query)
	return p
}

// Distinct 启用 DISTINCT 查询（去重）
func (p *PaginationBuilder) Distinct() *PaginationBuilder {
	p.distinct = true
	return p
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

	// 构建基础查询
	query := p.db.Model(p.model)

	// 应用连接查询
	for _, join := range p.joins {
		query = query.Joins(join)
	}

	// 应用查询条件
	for _, condition := range p.conditions {
		query = query.Where(condition.Query, condition.Args...)
	}

	// 应用分组
	for _, group := range p.groups {
		query = query.Group(group)
	}

	// 应用HAVING条件
	for _, having := range p.havings {
		query = query.Having(having.Query, having.Args...)
	}

	// 计算总记录数
	// 存在JOIN或多表查询时，用 WHERE id IN (子查询) 避免重复行 + 避免 Distinct().Count()/子查询表名引用 在 PostgreSQL 下的兼容问题
	if p.needDistinctCount() {
		tableName := p.resolveTableName()
		subQuery := p.db.Model(p.model).Select(tableName + ".id")
		for _, join := range p.joins {
			subQuery = subQuery.Joins(join)
		}
		for _, condition := range p.conditions {
			subQuery = subQuery.Where(condition.Query, condition.Args...)
		}
		for _, group := range p.groups {
			subQuery = subQuery.Group(group)
		}
		for _, having := range p.havings {
			subQuery = subQuery.Having(having.Query, having.Args...)
		}
		if err := p.db.Model(p.model).Where(tableName+".id IN (?)", subQuery).Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}
	}

	// 构建最终查询
	dataQuery := p.db.Model(p.model)

	// 应用SELECT字段
	if len(p.selects) > 0 {
		dataQuery = dataQuery.Select(strings.Join(p.selects, ","))
	}

	// 应用 DISTINCT
	if p.distinct {
		dataQuery = dataQuery.Distinct()
	}

	// 应用连接查询
	for _, join := range p.joins {
		dataQuery = dataQuery.Joins(join)
	}

	// 应用查询条件
	for _, condition := range p.conditions {
		dataQuery = dataQuery.Where(condition.Query, condition.Args...)
	}

	// 应用分组
	for _, group := range p.groups {
		dataQuery = dataQuery.Group(group)
	}

	// 应用HAVING条件
	for _, having := range p.havings {
		dataQuery = dataQuery.Having(having.Query, having.Args...)
	}

	// 应用排序
	for _, field := range p.orderFields {
		dataQuery = dataQuery.Order(field)
	}

	// 应用分页
	offset := (p.page - 1) * p.pageSize
	dataQuery = dataQuery.Offset(offset).Limit(p.pageSize)

	// 应用预加载（在分页后执行）
	for _, preload := range p.preloads {
		dataQuery = dataQuery.Preload(preload)
	}

	// 执行查询
	if err := dataQuery.Find(result).Error; err != nil {
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

// needDistinctCount 判断是否需要使用 DISTINCT 进行计数
// 当存在 JOIN 查询时，可能产生重复行，需要特殊处理
func (p *PaginationBuilder) needDistinctCount() bool {
	return len(p.joins) > 0 || len(p.groups) > 0
}

// resolveTableName 解析模型对应的数据库表名
// 返回带双引号的表名，以兼容 PostgreSQL 保留字（如 "user"）
func (p *PaginationBuilder) resolveTableName() string {
	switch m := p.model.(type) {
	case string:
		return m
	default:
		stmt := &gorm.Statement{DB: p.db}
		if err := stmt.Parse(p.model); err != nil {
			return ""
		}
		return `"` + stmt.Schema.Table + `"`
	}
}
