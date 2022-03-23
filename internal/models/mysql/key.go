package mysql

import (
	"kms/internal/models"
	"kms/internal/rule"

	"gorm.io/gorm"
)

// // TODO: remove
// var (
// 	v models.KeyRepo = (*KeyRepo)(nil)
// )

// KeyRepo KeyRepo
type KeyRepo struct {
}

// NewKeyRepo new
func NewKeyRepo() *KeyRepo {
	return &KeyRepo{}
}

// TableName tableName
func (k *KeyRepo) TableName() string {
	return models.KeyTableName
}

// Create create
func (k *KeyRepo) Create(db *gorm.DB, key *models.Key) error {
	return db.Model(key).Create(key).Error
}

// Delete delete
func (k *KeyRepo) Delete(db *gorm.DB, keyID string) int {
	return int(db.Table(k.TableName()).Where("key_id=? and active=?", keyID, rule.ActiveDisable).Delete(&models.Key{}).RowsAffected)
}

// Query query
func (k *KeyRepo) Query(db *gorm.DB, id string) (*models.Key, error) {
	e := &models.Key{}
	err := db.Table(k.TableName()).Where("key_id = ?", id).Find(e).Error
	if err != nil {
		return nil, err
	}

	return e, nil
}

// List list
func (k *KeyRepo) List(db *gorm.DB, userID string) ([]*models.SimplifyKey, int, error) {
	arr := make([]*models.SimplifyKey, 0)
	var total int64

	db = db.Table(k.TableName())

	if userID != "" {
		db = db.Where("owner = ?", userID)
	}

	db.Count(&total)
	err := db.Order("create_at DESC").Find(&arr).Error

	return arr, int(total), err
}

// Count count
func (k *KeyRepo) Count(db *gorm.DB, userID string) (int, error) {
	var total int64
	err := db.Table(k.TableName()).Where("owner = ?", userID).Count(&total).Error
	return int(total), err
}

// Active update status
func (k *KeyRepo) Active(db *gorm.DB, key *models.Key) error {
	return db.Table(k.TableName()).
		Where("key_id = ?", key.KeyID).
		Update("active", key.Active).Error
}

// Update update
func (k *KeyRepo) Update(db *gorm.DB, key *models.Key) error {
	mp := map[string]interface{}{
		"title":       key.Title,
		"description": key.Description,
		"assignee":    key.Assignee,
	}
	return db.Table(k.TableName()).Where("key_id=?", key.KeyID).Updates(mp).Error
}
