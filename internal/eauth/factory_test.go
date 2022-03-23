package eauth

import (
	"fmt"
	"kms/internal/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func _TestFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}

func (fs *FactorySuite) SetupSuite() {

}

func (fs *FactorySuite) TestCreate() {
	ek := &models.AgencyKey{
		AuthType:    "signature",
		KeyID:       "UjjNR4B_qO6925q7afx8dQ",
		KeySecret:   "FgePLETu5DMLh4yvh_0h7Gbd-Tv9Q2qZKI2hYBT4OrUiC7ONaS4RNC46AygVVQnNF8M8qNUL9Vytgo0",
		AuthContent: `{"cmds": "sort query gonic asc|append begin GET\n/saas/\n|sha256 <SECRET_KEY>|base64 std encode", "in": "body", "name": "name"}`,
	}

	cases := []*models.AgencyKey{
		{
			AuthType: "dsadsa",
		},
		{
			AuthType:    "signature",
			KeyID:       "UjjNR4B_qO6925q7afx8dQ",
			KeySecret:   "FgePLETu5DMLh4yvh_0h7Gbd-Tv9Q2qZKI2hYBT4OrUiC7ONaS4RNC46AygVVQnNF8M8qNUL9Vytgo0",
			AuthContent: `{"cmds": "sort query gonic asc|append begin GET\n/saas/\n|sha256 <SECRET_KEY>|base64 std encode", "in": "body", "name": "name"}`,
		},
	}

	for _, c := range cases {
		auth, err := factory.Create(c)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		token, _ := auth.Invoke(ek)
		fmt.Printf("token: %#v\n", token)
	}

}
