package model

import (
	"time"
)

type Transaction struct {
	ID          int        `gorm:"primary_key;AUTO_INCREMENT" json:"inspection_id"`
	FromAddress string     `gorm:"not null" json:"from_address"`
	ToAddress   string     `gorm:"not null" json:"to_address"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:now()" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (Transaction) TableName() string {
	return "transaction"
}
