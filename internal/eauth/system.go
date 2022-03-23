package eauth

import (
	"kms/internal/enums"
	"kms/internal/models"
	"kms/pkg/signature"
)

// System system auth
type System struct {
	AuthBase
}

// SystemContent system conten
type SystemContent struct {
	emptyAuthContentParser
}

// Name name
func (s *System) Name() enums.Enum {
	return enums.AuthSystem
}

// Invoke invoke
func (s *System) Invoke(d interface{}) ([]*AuthResp, error) {
	sign, err := signature.Signature(d.(map[string]interface{}), s.secret)
	return []*AuthResp{{
		Name:  "system",
		In:    "system",
		Value: sign,
	}}, err
}

// Init init
func (s *System) Init(ak *models.AgencyKey) error {
	s.authContent = &SystemContent{}
	return s.AuthBase.Init(ak)
}
