package types

type RoleMenuIdsRequest struct {
	RoleID int64 `json:"role_id" form:"role_id" binding:"required"`
}

type RoleMenuRequest struct {
	RoleID int64 `json:"role_id" form:"role_id" binding:"required"`
}

type AddRoleMenuRequest struct {
	RoleID  int64   `json:"role_id" binding:"required"`
	MenuIDs []int64 `json:"menu_ids" binding:"required"`
}
