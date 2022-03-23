package eauth

import (
	"encoding/json"
	"fmt"
	"kms/internal/models"
	"kms/internal/rule"
	"kms/pkg/crypto/encode"
	"testing"
)

func TestSignature(t *testing.T) {
	s := Signature{}
	k := &models.AgencyKey{
		KeySecret: "kkkkSecret2",
		KeyID:     "kkkkID",
		AuthContent: `
[
	{
		"type":"signcmd",
		"data":"sort query gonic asc|append begin GET\n/saas/\n|sha256 <SECRET_KEY>|base64 std encode"
	},
	{
		"name":"access_key_id",
		"type":"keyid",
		"in":"query"
	},
	{
		"name":"Signature",
		"type":"signature",
		"in":"header"	
	}
]
`,
	}
	e, err := encode.SecretEncodeString(k.KeySecret, k.KeyID)
	if err != nil {
		panic(err)
	}
	k.KeySecret = e
	if _, err := rule.ParseAuthContent(k.AuthContent, true); err != nil {
		panic(err)
	}
	if err := s.Init(k); err != nil {
		panic(err)
	}
	resp, err := s.Invoke("{}")
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
