package eauth

import (
	"kms/internal/enums"
	"kms/internal/models"
	"kms/internal/rule"
	"kms/pkg/security/security"
	"strings"

	error2 "github.com/quanxiang-cloud/cabin/error"
)

// Signature signature
type Signature struct {
	AuthBase
}

// SignContent auth content of sign
type SignContent struct {
	Cmds string   `json:"cmds"`
	Key  AuthResp `json:"key"`
	Sign AuthResp `json:"sign"`
}

// Parse parse and verify config values for this auth type
func (c *SignContent) Parse(vals rule.ConfigValueSet) error {
	for _, v := range vals {
		switch {
		case v.In == "" && v.Type == enums.BasicTypeSignCmd.Val():
			c.Cmds = v.Data
		case v.Type == enums.BasicTypeKeyID.Val():
			c.Key.Name = v.Name
			c.Key.In = v.In
		case v.Type == enums.BasicTypeSignature.Val():
			c.Sign.Name = v.Name
			c.Sign.In = v.In
		}
	}
	if c.Cmds == "" {
		return error2.NewErrorWithString(error2.ErrParams,
			"missing signcmd config")
	}
	if c.Key.Name == "" {
		return error2.NewErrorWithString(error2.ErrParams,
			"missing keyid response config")
	}
	if c.Sign.Name == "" {
		return error2.NewErrorWithString(error2.ErrParams,
			"missing signature response config")
	}
	return nil
}

// Name name
func (s *Signature) Name() enums.Enum {
	return enums.AuthSignature
}

// Init overwrite init func, instantiate authContent then call super init
func (s *Signature) Init(ak *models.AgencyKey) error {
	s.authContent = &SignContent{}
	return s.AuthBase.Init(ak)
}

// Invoke to signature
func (s *Signature) Invoke(d interface{}) ([]*AuthResp, error) {
	c := s.authContent.(*SignContent)
	cmdRun := strings.ReplaceAll(c.Cmds, security.ConfigSecretKeyName, s.secret)
	//BUG: should append access-key-id here
	if p := &c.Key; p.Name != "" {
		if m, ok := d.(map[string]interface{}); ok {
			m[p.Name] = s.keyID
		}
	}
	sign, _ := security.Signature(cmdRun, d)
	resp := make([]*AuthResp, 0, 2)
	if p := &c.Key; p.Name != "" {
		resp = append(resp, &AuthResp{
			Name:  p.Name,
			In:    p.In,
			Value: s.keyID,
		})
	}
	if p := &c.Sign; p.Name != "" {
		resp = append(resp, &AuthResp{
			Name:  p.Name,
			In:    p.In,
			Value: sign,
		})
	}
	return resp, nil
}
