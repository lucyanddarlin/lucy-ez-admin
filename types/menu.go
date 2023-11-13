package types

type AddMenuRequest struct {
	ParentID   int64  `json:"parent_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Type       string `json:"type" binding:"required"`
	Permission string `json:"permission"`
	Method     string `json:"method"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Weight     string `json:"weight"`
	IsHIdden   bool   `json:"is_hidden"`
	IsCache    bool   `json:"is_cache"`
	IsHome     bool   `json:"is_home"`
}
