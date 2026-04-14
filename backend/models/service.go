package models

import (
	"time"
)

// Service 服务模型
type Service struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	GitUrl      string    `gorm:"size:500" json:"gitUrl"`
	Category    string    `gorm:"size:50" json:"category"`
	SubCategory string    `gorm:"size:50" json:"subCategory"`
	Remark      string    `gorm:"size:500" json:"remark"`
	// ServerID 默认部署目标机器；创建服务时必填，自动部署优先使用
	ServerID uint `json:"serverId" gorm:"column:server_id;index;comment:关联服务器ID"`
}
