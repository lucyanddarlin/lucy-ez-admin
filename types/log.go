package types

type LoginLogRequest struct {
	Page     int    `json:"page" form:"page" binding:"required" sql:"-"`
	PageSize int    `json:"page_size" form:"page_size" binding:"required,max=50" sql:"-"`
	Phone    string `json:"phone" form:"phone"`
	Status   *bool  `json:"status" form:"status"`
	Start    int64  `json:"start" form:"start" sql:"> ?" column:"create_at"`
	End      int64  `json:"end" form:"end" sql:"< ?" column:"create_at"`
}
