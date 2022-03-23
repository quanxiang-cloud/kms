package models

import (
	"gorm.io/gorm"
)

// AgencyKey external key
type AgencyKey struct {
	ID          string //unique id
	Owner       string //owner id
	OwnerName   string
	Name        string
	Title       string
	Description string
	Service     string //belong service, eg: system_form
	Host        string //service host, eg: api.xxx.com:8080
	AuthType    string //signature/cookie/oauth2...
	AuthContent string //Authorize detail
	KeyID       string //key id
	KeySecret   string //crypt key secret
	KeyContent  string //key content
	Parsed      int    //1 parsed 0 not parse
	Active      int    //1 active 0 disable
	CreateAt    int64  //create time
	UpdateAt    int64  //update time
}

// SimplifyAgencyKey key info except secret
type SimplifyAgencyKey struct {
	ID          string `json:"id"`    //unique id
	Owner       string `json:"owner"` //owner id
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`     //belong service, eg: system_form
	Host        string `json:"host"`        //service host, eg: api.xxx.com:8080
	AuthType    string `json:"authType"`    //signature/cookie/oauth2...
	AuthContent string `json:"authContent"` //Authorize detail
	KeyID       string `json:"keyID"`       //key id
	KeyContent  string `json:"keyContent"`  //key content
	Active      int    `json:"active"`      //1 active 0 disable
	CreateAt    int64  `json:"createAt"`    //create time
	UpdateAt    int64  `json:"updateAt"`    //update time
}

// AgencyKeyTableName ext_key table name
const AgencyKeyTableName = "customer_secret_key"

// TableName table name
func (e *AgencyKey) TableName() string {
	return AgencyKeyTableName
}

// AgencyKeyRepo repo
type AgencyKeyRepo interface {
	Create(db *gorm.DB, ak *AgencyKey) error
	Delete(db *gorm.DB, keyID string) error
	DeleteInBatch(db *gorm.DB, service []string) ([]string, error)
	Query(db *gorm.DB, keyID string) (*AgencyKey, error)
	QueryInService(db *gorm.DB, service, keyID string) (*AgencyKey, error)
	List(db *gorm.DB, service string) ([]*SimplifyAgencyKey, int, error)
	Update(db *gorm.DB, id, title, description string) error
	UpdateContent(db *gorm.DB, id, authContent string, parsed int) error
	UpdateInBatch(db *gorm.DB, host, service, authType, authContent string) ([]string, error)
	Active(db *gorm.DB, keyID string, status int) error
	DelByPrefixPath(db *gorm.DB, path string) ([]*AgencyKey, error)
}
