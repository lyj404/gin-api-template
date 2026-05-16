package dto

type CreateOrgUnitRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint64  `json:"parent_id,string"`
}

type UpdateOrgUnitRequest struct {
	Name     string `json:"name"`
	ParentID *uint64  `json:"parent_id,string"`
}

type OrgUnitResponse struct {
	ID       uint64   `json:"id,string"`
	Name     string `json:"name"`
	ParentID *uint64  `json:"parent_id,string"`
	Path     string `json:"path"`
	Level    int    `json:"level"`
}
