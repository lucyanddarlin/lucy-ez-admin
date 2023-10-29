package types

type BaseModel struct {
	ID        int64 `json:"id" gorm:"primary_key;autoIncrement;size:32;comment:主键ID"`
	CreatedAt int64 `json:"created_at,omitempty" gorm:"index;comment:创建时间"`
	UpdatedAt int64 `json:"updated_at,omitempty" gorm:"index;comment:修改时间"`
}
