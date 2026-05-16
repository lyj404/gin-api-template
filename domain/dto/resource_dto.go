package dto

type BindResourceRequest struct {
	RoleID  uint64 `json:"role_id,string" binding:"required"`
	IsWrite bool `json:"is_write"`
}

type CreateResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Pattern     string `json:"pattern" binding:"required"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type UpdateResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Pattern     string `json:"pattern" binding:"required"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type ResourceListResponse struct {
	ID          uint64   `json:"id,string"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Method      string `json:"method"`
	Entity      string `json:"entity"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
