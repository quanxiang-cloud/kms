package eauth

import (
	"kms/internal/enums"
	"kms/internal/models"
)

// None none
type None struct {
	AuthBase
}

// NoneContent noneContent
type NoneContent struct {
	emptyAuthContentParser
}

// Name name
func (n *None) Name() enums.Enum {
	return enums.AuthNone
}

// Invoke invoke
func (n *None) Invoke(d interface{}) ([]*AuthResp, error) {
	return []*AuthResp{}, nil
}

// Init init
func (n *None) Init(ek *models.AgencyKey) error {
	n.authContent = &NoneContent{}
	return n.AuthBase.Init(ek)
}
