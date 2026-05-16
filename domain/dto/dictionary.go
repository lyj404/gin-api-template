package dto

// CreateDictRequest 创建字典请求
type CreateDictRequest struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}

// UpdateDictRequest 更新字典请求
type UpdateDictRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}

// DictResponse 字典响应
type DictResponse struct {
	ID      uint64              `json:"id"`
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Status  int               `json:"status"`
	Desc    string            `json:"desc"`
	Details []DictDetailResponse `json:"details,omitempty"`
}

// CreateDictDetailRequest 创建字典详情请求
type CreateDictDetailRequest struct {
	DictID uint64 `json:"dict_id" binding:"required"`
	Label  string `json:"label" binding:"required"`
	Value  string `json:"value" binding:"required"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// UpdateDictDetailRequest 更新字典详情请求
type UpdateDictDetailRequest struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// DictDetailResponse 字典详情响应
type DictDetailResponse struct {
	ID     uint64 `json:"id"`
	DictID string `json:"dict_id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}
