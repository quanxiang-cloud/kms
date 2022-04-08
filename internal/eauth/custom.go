package eauth

import (
	"kms/internal/enums"
	"kms/internal/models"
)

// Custom abstract struct of custom defination
type Custom struct {
	name enums.Enum
	AuthBase
}

// SetName set method name
func (c *Custom) SetName(name string) {
	c.name = enums.Enum(name)
}

// Name name
func (c *Custom) Name() enums.Enum {
	return c.name
}

// Init init
func (c *Custom) Init(ak *models.AgencyKey) error {
	return c.AuthBase.Init(ak)
}

// Invoke invoke
func (c *Custom) Invoke(d interface{}) ([]*AuthResp, error) {
	// TODO:
}
