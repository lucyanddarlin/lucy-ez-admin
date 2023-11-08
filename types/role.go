package types

type AddRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Keyword     string `json:"keyword" binding:"required"`
	Status      *bool  `json:"status" binding:"required"`
	Weight      int    `json:"weight"`
	ParentID    int64  `json:"parent_id"`
	Description string `json:"description"`
	TeamIds     string `json:"team_ids"`
	DataScope   string `json:"data_scope" binding:"required"`
}

type UpdateRoleRequest struct {
	ID          int64   `json:"id" binding:"required"`
	Name        string  `json:"name"`
	Status      *bool   `json:"status"`
	Weight      int     `json:"weight"`
	Description string  `json:"description"`
	DataScope   string  `json:"dataScope"`
	TeamIds     *string `json:"team_id"`
	ParentID    int64   `json:"parent_id"`
}

type DeleteRoleRequest struct {
	ID int64 `json:"id" binding:"required"`
}
