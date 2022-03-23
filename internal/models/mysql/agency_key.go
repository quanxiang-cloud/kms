package mysql

import (
	"fmt"
	"kms/internal/models"
	"kms/internal/rule"

	"gorm.io/gorm"
)

var (
	a models.AgencyKeyRepo = (*AgencyKeyRepo)(nil)
)

// AgencyKeyRepo ext_key repo
type AgencyKeyRepo struct {
}

// CreateExtKeyRepo create ext_key repo
func CreateExtKeyRepo() *AgencyKeyRepo {
	return &AgencyKeyRepo{}
}

// TableName table name
func (e *AgencyKeyRepo) TableName() string {
	return models.AgencyKeyTableName
}

// Create create
func (e *AgencyKeyRepo) Create(db *gorm.DB, key *models.AgencyKey) error {
	return db.Model(key).Create(key).Error
}

// Delete delete
func (e *AgencyKeyRepo) Delete(db *gorm.DB, id string) error {
	return db.Table(e.TableName()).Where("id = ?", id).Delete(&models.AgencyKey{}).Error
}

// DeleteInBatch delete in batch
func (e *AgencyKeyRepo) DeleteInBatch(db *gorm.DB, service []string) ([]string, error) {
	sql := db.Table(e.TableName()).Where("service in (?)", service)

	entities := make([]*models.AgencyKey, 0)
	err := sql.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(entities))
	for _, v := range entities {
		ids = append(ids, v.ID)
	}
	err = sql.Delete(&models.AgencyKey{}).Error
	return ids, err
}

// Query query
func (e *AgencyKeyRepo) Query(db *gorm.DB, id string) (*models.AgencyKey, error) {
	ek := &models.AgencyKey{}
	err := db.Table(e.TableName()).Where("id = ?", id).Find(ek).Error
	if err != nil {
		return nil, err
	}
	return ek, nil
}

// QueryInService QueryInService
func (e *AgencyKeyRepo) QueryInService(db *gorm.DB, service, keyID string) (*models.AgencyKey, error) {
	ret := &models.AgencyKey{}
	err := db.Table(e.TableName()).Where("service=? and key_id=?", service, keyID).Find(ret).Error
	return ret, err
}

// List list
func (e *AgencyKeyRepo) List(db *gorm.DB, service string) ([]*models.SimplifyAgencyKey, int, error) {
	arr := make([]*models.SimplifyAgencyKey, 0)
	var total int64

	db = db.Table(e.TableName())

	if service != "" {
		db = db.Where("service = ?", service)
	}

	db.Count(&total)
	err := db.Order("create_at DESC").Find(&arr).Error

	return arr, int(total), err
}

// Update update
func (e *AgencyKeyRepo) Update(db *gorm.DB, id, title, description string) error {
	return db.Table(e.TableName()).Where("id=?", id).Updates(map[string]interface{}{
		"title":       title,
		"description": description,
	}).Error
}

// UpdateContent update auth content
func (e *AgencyKeyRepo) UpdateContent(db *gorm.DB, id, authContent string, parsed int) error {
	return db.Table(e.TableName()).Where("id=?", id).Updates(map[string]interface{}{
		"auth_content": authContent,
		"parsed":       parsed,
	}).Error
}

// UpdateInBatch update authType & authContent by service
func (e *AgencyKeyRepo) UpdateInBatch(db *gorm.DB, host, service, authType, authContent string) ([]string, error) {
	keys := make([]*models.AgencyKey, 0)
	updateDB := db.Table(e.TableName()).Where("service=?", service)
	err := updateDB.Find(&keys).Error
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(keys))
	if len(keys) < 1 {
		return ids, nil
	}
	for k, v := range keys {
		ids[k] = v.ID
	}

	mp := map[string]interface{}{
		"host":         host,
		"auth_type":    authType,
		"auth_content": authContent,
		"parsed":       rule.NotParsed,
	}
	err = updateDB.Updates(mp).Error
	return ids, err
}

// Active update status
func (e *AgencyKeyRepo) Active(db *gorm.DB, keyID string, active int) error {
	return db.Table(e.TableName()).Where("id=?", keyID).Update("active", active).Error
}

// DelByPrefixPath DelByPrefixPath
func (e *AgencyKeyRepo) DelByPrefixPath(db *gorm.DB, path string) ([]*models.AgencyKey, error) {
	ret := make([]*models.AgencyKey, 0)
	del := db.Table(e.TableName()).Where("service like ?", fmt.Sprintf("%s%%", path))
	err := del.Find(&ret).Error
	if err != nil {
		return nil, err
	}
	err = del.Delete(&models.AgencyKey{}).Error
	return ret, err
}
