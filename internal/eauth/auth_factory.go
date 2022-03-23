package eauth

import (
	"fmt"
	"kms/internal/models"
	"kms/internal/rule"
	"reflect"
)

func init() {
	factory.MustReg(&None{
		AuthBase: AuthBase{
			authContentVals: rule.ConfigValueSet{},
		},
	})
	factory.MustReg(&System{
		AuthBase: AuthBase{
			authContentVals: rule.ConfigValueSet{},
		},
	})
	factory.MustReg(&Signature{
		AuthBase: AuthBase{
			authContentVals: rule.ConfigValueSet{
				{
					Name: "access_key_id",
					Type: "keyid",
					In:   "query",
				},
				{
					Name: "signature",
					Type: "signature",
					In:   "body",
				},
				{
					Type: "signcmd",
					Data: "sort query gonic asc|append begin GET\n/iaas/\n|sha256 <SECRET_KEY>|base64 std encode",
				},
			},
		},
	})
	factory.MustReg(&Cookie{
		AuthBase: AuthBase{
			authContentVals: rule.ConfigValueSet{
				{
					Name: "userName",
					Type: "keyid",
					From: "query",
				},
				{
					Name: "password",
					Type: "keysecret",
					From: "query",
				},
				{
					Name: "serviceName",
					Type: "string",
					Data: "clogin",
					From: "query",
				},
				{
					Type: "authurl",
					Data: "http://host/api",
				},
				{
					Type: "method",
					Data: "GET",
				},
				{
					Name:      "binding",
					Type:      "string",
					CaptureIn: "body",
					In:        "query",
				},
				{
					Type: "expire",
					Data: "5",
				},
			},
		},
	})
}

var (
	factory = newAuthFactory()
)

// GetAuthFactory get
func GetAuthFactory() *AuthFactory {
	return factory
}

// AuthFactory auth type factory
type AuthFactory struct {
	mp     map[string]reflect.Type
	sample map[string]interface{}
}

// newAuthFactory instant factory, only in this package
func newAuthFactory() *AuthFactory {
	return &AuthFactory{
		mp:     make(map[string]reflect.Type),
		sample: make(map[string]interface{}),
	}
}

// MustReg register auth component, panic err if type is multiple
func (f *AuthFactory) MustReg(c Auth) {
	err := f.Reg(c)
	if err != nil {
		panic(err)
	}
}

// Reg register auth component
func (f *AuthFactory) Reg(c Auth) error {
	if _, ok := f.mp[c.Name().Val()]; ok {
		return fmt.Errorf("duplicate reg of %s, %#v", c.Name(), c)
	}
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	f.mp[c.Name().Val()] = t
	f.sample[c.Name().Val()] = c
	return nil
}

// Create create auth
func (f *AuthFactory) Create(d *models.AgencyKey) (Auth, error) {
	t, ok := f.mp[d.AuthType]
	if ok {
		rt := reflect.New(t)
		inst := rt.Interface().(Auth)
		if err := inst.Init(d); err != nil {
			return nil, err
		}
		return inst, nil
	}
	err := fmt.Errorf("auth factory can't create %s", d.AuthType)
	return nil, err
}

// CreateSample create sample
func (f *AuthFactory) CreateSample(name string) (interface{}, error) {
	sample, ok := f.sample[name]
	if !ok {
		return nil, fmt.Errorf("not registed type: %s", name)
	}
	return sample, nil
}
