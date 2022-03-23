package rule

import (
	"encoding/json"
	"fmt"
	"kms/internal/enums"
	"kms/pkg/security/security"
	"strconv"

	error2 "github.com/quanxiang-cloud/cabin/error"
)

const (
	req  = "request"
	resp = "response"
)

// ConfigValue is the value defined for auth content
type ConfigValue struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Data      string `json:"data,omitempty"`
	In        string `json:"in,omitempty"`        // empty means don't response
	From      string `json:"from,omitempty"`      // describe http params
	CaptureIn string `json:"captureIn,omitempty"` // reserve for oauth2
}

// Validate verify the value
func (c ConfigValue) Validate() error {
	if c.In == "" && c.From == "" {
		if c.Data == "" {
			return error2.NewErrorWithString(error2.ErrParams,
				fmt.Sprintf("missing data for value %s", c.Name))
		}
		if !enums.ValueTypeEnumSet.Verify(c.Type) {
			return error2.NewErrorWithString(error2.ErrParams,
				fmt.Sprintf("input value %s type(%s) outof %v", c.Name, c.Type, enums.ValueTypeEnumSet.GetAll()))
		}
	}
	if c.In != "" && !enums.ValueInEnumSet.Verify(c.In) {
		return error2.NewErrorWithString(error2.ErrParams,
			fmt.Sprintf("input value %s in(%s) outof %v", c.Name, c.In, enums.ValueInEnumSet.GetAll()))
	}
	if c.From != "" && !enums.ValueFromEnumSet.Verify(c.From) {
		return error2.NewErrorWithString(error2.ErrParams,
			fmt.Sprintf("input value %s in(%s) outof %v", c.Name, c.In, enums.ValueFromEnumSet.GetAll()))
	}
	if !enums.BasicTypeEnumSet.Verify(c.Type) {
		return error2.NewErrorWithString(error2.ErrParams,
			fmt.Sprintf("input value %s type(%s) outof %v", c.Name, c.Type, enums.BasicTypeEnumSet.GetAll()))
	}
	switch c.Type {
	case enums.BasicTypeNumber.Val(), enums.BasicTypeExpire.Val(), enums.BasicTypeRandom.Val():
		if _, err := strconv.ParseFloat(c.Data, 64); err != nil {
			return error2.NewErrorWithString(error2.ErrParams, fmt.Sprintf("input value %s type(%s) is invalid", c.Name, c.Type))
		}
	case enums.BasicTypeBoolean.Val():
		if _, err := strconv.ParseBool(c.Data); err != nil {
			return error2.NewErrorWithString(error2.ErrParams, fmt.Sprintf("input value %s type(%s) is invalid", c.Name, c.Type))
		}
	case enums.BasicTypeMethod.Val():
		if !enums.HTTPMethodSet.Verify(c.Data) {
			return fmt.Errorf("method only in %v", enums.HTTPMethodSet.GetAll())
		}
	case enums.BasicTypeSignCmd.Val():
		cmds := security.SplitSignatureCMD(c.Data)
		if len(cmds) == 0 {
			return fmt.Errorf("no cmd found")
		}
		for _, v := range cmds {
			if err := security.ValidSignatureCmds(v.CMD); err != nil {
				return err
			}
		}
		// TODO: check
	default:
		// do nothing
	}
	return nil
}

// ConfigValueSet is value set
type ConfigValueSet []ConfigValue

// Validate verify the value set
func (s ConfigValueSet) Validate() error {
	for _, v := range s {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ParseAuthContent parse and validate authCOntent
func ParseAuthContent(authContent string, validate bool) (ConfigValueSet, error) {
	var d ConfigValueSet
	if err := json.Unmarshal([]byte(authContent), &d); err != nil {
		return nil, err
	}
	if validate {
		if err := d.Validate(); err != nil {
			return nil, err
		}
	}
	return d, nil
}
