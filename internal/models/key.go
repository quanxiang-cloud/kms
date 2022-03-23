package models

import (
	"gorm.io/gorm"
)

// Key key
type Key struct {
	ID          string `gorm:"primarykey"`
	Owner       string
	OwnerName   string
	Name        string
	Title       string
	Description string
	KeyID       string
	KeySecret   string
	Active      int
	Assignee    string
	CreateAt    int64
	UpdateAt    int64
}

// SimplifyKey simplify key
type SimplifyKey struct {
	ID          string `gorm:"primarykey" json:"id"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	KeyID       string `json:"keyID"`
	Active      int    `json:"active"`
	Assignee    string `json:"assignee"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

// KeyTableName table name in mysql
const KeyTableName = "secret_key"

// TableName table name of Key struct
func (k *Key) TableName() string {
	return KeyTableName
}

// KeyRepo KeyRepo
type KeyRepo interface {
	Create(db *gorm.DB, key *Key) error
	Delete(db *gorm.DB, keyID string) int
	Query(db *gorm.DB, id string) (*Key, error)
	List(db *gorm.DB, userID string) ([]*SimplifyKey, int, error)
	Count(db *gorm.DB, userID string) (int, error)
	Update(db *gorm.DB, key *Key) error
	Active(db *gorm.DB, key *Key) error
}
