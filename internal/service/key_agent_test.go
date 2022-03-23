package service

import (
	"context"
	"fmt"
	"kms/pkg/misc/config"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExtKeySuite struct {
	suite.Suite

	ctx    context.Context
	conf   *config.Config
	r      *gin.Engine
	extKey KeyAgent

	key string
}

func _TestExtKeySuite(t *testing.T) {
	suite.Run(t, &ExtKeySuite{})
}

func (suite *ExtKeySuite) SetupSuite() {
	var err error
	suite.conf, err = config.NewConfig("./testdata/test/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)
	assert.Nil(suite.T(), err)
}

func (suite *ExtKeySuite) TestCreate() {
	fmt.Println("=================================Test Create=================================")
	case1 := &UploadAgencyReq{
		UserID:      "c99b9fbc-9a4e-43c0-b3e2-9c5dc55782a0",
		Name:        "from test: name",
		Title:       "from test: title",
		Description: "from test: description",
		Service:     "from test: service",
		Host:        "from test: host",
		AuthType:    "signature",
		AuthContent: `{"cmd": "from test"}`,
		KeyID:       "key id",
		KeySecret:   "key secret",
		KeyContent:  "",
	}
	resp, err := suite.extKey.Upload(suite.ctx, case1)
	assert.Nil(suite.T(), err)
	fmt.Println("id: ", resp.ID)
	fmt.Println("=================================    End    =================================")
}

func (suite *ExtKeySuite) TestQuery() {
	fmt.Println("=================================Test  Query=================================")
	case1 := &QueryAgencyReq{
		ID: "",
	}
	resp, err := suite.extKey.Query(suite.ctx, case1)
	assert.Nil(suite.T(), err)
	fmt.Printf("%+v\n", resp)
	fmt.Println("=================================    End    =================================")
}

func (suite *ExtKeySuite) TestList() {

}

func (suite *ExtKeySuite) TestDel() {
	fmt.Println("=================================Test  Query=================================")
	case1 := &DeleteAgencyReq{
		ID: "d0265fe3-95ce-443b-b193-f188cffc7477",
	}
	resp, err := suite.extKey.Delete(suite.ctx, case1)
	assert.Nil(suite.T(), err)
	_ = resp
	fmt.Println("=================================    End    =================================")
}
