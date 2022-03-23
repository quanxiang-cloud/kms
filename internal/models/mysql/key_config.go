package mysql

import (
	"kms/internal/models"

	"gorm.io/gorm"
)

// KeyConfigRepo KeyConfigRepo
type KeyConfigRepo struct {
}

// NewKeyConfigRepo new
func NewKeyConfigRepo() *KeyConfigRepo {
	return &KeyConfigRepo{}
}

// TableName tableName
func (k *KeyConfigRepo) TableName() string {
	return models.KeyConfigTableName
}

// Create create
func (k *KeyConfigRepo) Create(db *gorm.DB, config *models.KeyConfig) error {
	return db.Model(config).Create(config).Error
}

// Update update
func (k *KeyConfigRepo) Update(db *gorm.DB, config *models.KeyConfig) error {
	mp := map[string]interface{}{
		"config_content": config.ConfigContent,
	}
	return db.Table(k.TableName()).Where("owner=?", config.Owner).Updates(mp).Error
}

// Query query
func (k *KeyConfigRepo) Query(db *gorm.DB, owner string) (*models.KeyConfig, error) {
	kc := &models.KeyConfig{}
	err := db.Table(k.TableName()).Where("owner=?", owner).Find(kc).Error
	return kc, err
}

// Delete delete
func (k *KeyConfigRepo) Delete(db *gorm.DB, owner string) error {
	return db.Table(k.TableName()).Where("owner=?", owner).Delete(&models.KeyConfig{}).Error
}
