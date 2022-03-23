package eauth

import (
	"encoding/json"
	"kms/internal/enums"
	"kms/internal/models"
	"kms/internal/models/redis"
	"kms/internal/rule"
	"kms/pkg/crypto/decode"
	"kms/pkg/dbcli"
	"reflect"
	"unsafe"
)

// Config
const (
	ConfigKeyName       = "${KEY_ID}"
	ConfigSecretKeyName = "${SECRET_KEY}"
	// ConfigTime          = "${TIME{2006.01.02 15:04:05.9999999 -0700}}"
	// ConfigRandom        = "${RANDOM{5}}"
)

// Auth auth
type Auth interface {
	Invoke(d interface{}) ([]*AuthResp, error)
	Name() enums.Enum
	Init(ak *models.AgencyKey) error
	GetContent() interface{}
	GetContentVals() rule.ConfigValueSet
}

// AuthResp resp
type AuthResp struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type emptyAuthContentParser struct{}

func (t emptyAuthContentParser) Parse(vals rule.ConfigValueSet) error {
	return nil
}

type authContentParser interface {
	Parse(vals rule.ConfigValueSet) error
}

// AuthBase base
type AuthBase struct {
	agencyKey       *models.AgencyKey
	keyID           string
	secret          string
	authContent     authContentParser
	authContentVals rule.ConfigValueSet
}

// Init auth
// deserialize authContent if it was instantiated by child init func
func (b *AuthBase) Init(ak *models.AgencyKey) error {
	b.agencyKey, b.keyID = ak, ak.KeyID
	secret, err := decode.SecretDecodeString(b.agencyKey.KeySecret, b.agencyKey.KeyID)
	if err != nil {
		return err
	}
	b.secret = secret

	// defualt init authContent
	if ak.AuthContent != "" {
		if err := b.initAuthContent(); err != nil {
			return err
		}
	}
	return nil
}

func (b *AuthBase) initAuthContent() error {
	if !rule.CheckParse(b.agencyKey.Parsed) {
		if err := json.Unmarshal(unsafeString2Bytes(b.agencyKey.AuthContent), &b.authContentVals); err != nil {
			return err
		}
		if b.authContent != nil {
			if err := b.authContent.Parse(b.authContentVals); err != nil {
				return err
			}
		}
	} else {
		if err := json.Unmarshal([]byte(b.agencyKey.AuthContent), &b.authContent); err != nil {
			return err
		}
	}
	return nil
}

// GetContent get auth content
func (b *AuthBase) GetContent() interface{} {
	return b.authContent
}

// GetContentVals get auth content configs
func (b *AuthBase) GetContentVals() rule.ConfigValueSet {
	return b.authContentVals
}

// unsafeString2Bytes convert string to []byte without copy
func unsafeString2Bytes(src string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&src))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func getRedisCli() (*redis.Client, error) {
	c, err := dbcli.GetRedisClient(nil)
	if err != nil {
		return nil, err
	}
	return redis.NewRedisClient(c), nil
}
