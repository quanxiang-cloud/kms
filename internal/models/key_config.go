package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

var (
	// DefaultNum key default num, it will be inited
	DefaultNum = -1
)

// KeyConfig key conf
type KeyConfig struct {
	ID            string
	Owner         string
	OwnerName     string
	ConfigContent ConfigContent
	CreateAt      int64
	UpdateAt      int64
}

// ConfigContent configContents
type ConfigContent map[string]string

// KeyConfigTableName KeyConfigTableName
const KeyConfigTableName = "secret_key_config"

// TableName tableName
func (kc *KeyConfig) TableName() string {
	return KeyConfigTableName
}

// KeyConfigRepo key config repo
type KeyConfigRepo interface {
	Create(db *gorm.DB, conf *KeyConfig) error
	Update(db *gorm.DB, conf *KeyConfig) error
	Query(db *gorm.DB, owner string) (*KeyConfig, error)
	Delete(db *gorm.DB, owner string) error
}

// Value marshal
func (kc ConfigContent) Value() (driver.Value, error) {
	//return json.MarshalIndent(c, "", "  ")
	return json.Marshal(kc)
}

// Scan unmarshal
func (kc *ConfigContent) Scan(data interface{}) error {
	d, ok := data.([]byte)
	if !ok {
		return errors.New("unknown scan data type")
	}
	if err := json.Unmarshal(d, kc); err != nil {
		return err
	}
	return nil
}
