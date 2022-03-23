package service

import (
	"context"
	"fmt"
	"kms/internal/rule"
	"kms/pkg/misc/config"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KeyGenSuite struct {
	suite.Suite

	ctx    context.Context
	conf   *config.Config
	r      *gin.Engine
	keyGen KeyGenerator

	key string
}

func _TestKeyGen(t *testing.T) {
	suite.Run(t, new(KeyGenSuite))
}

func (suite *KeyGenSuite) SetupSuite() {
	var err error
	suite.conf, err = config.NewConfig("./testdata/test/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)

	log := logger.New(&suite.conf.Logger)

	suite.keyGen, err = CreateKeyGenerator(suite.conf, log)
	assert.Nil(suite.T(), err)
}

func (suite *KeyGenSuite) TestCreateKey() {
	req := &CreateReq{
		UserID:      "c99b9fbc-9a4e-43c0-b3e2-9c5dc55782a0",
		Name:        "test2",
		Title:       "title",
		Description: "description",
		Assignee:    "assignee",
	}

	resp, err := suite.keyGen.CreateKey(suite.ctx, req)
	assert.Nil(suite.T(), err)
	fmt.Println("=================Create===================")
	fmt.Printf("%+v\n", resp)
	suite.key = resp.KeyID
}

func (suite *KeyGenSuite) TestQueryKey() {
	req := &QueryReq{
		// KeyID: suite.key,
		KeyID: "KTdcf7HesCicFNo2ZP58jQ",
	}

	resp, err := suite.keyGen.QueryKey(suite.ctx, req)
	assert.Nil(suite.T(), err)
	fmt.Println("=================Query===================")
	fmt.Printf("%+v\n", resp)
}

func (suite *KeyGenSuite) TestUpdateKey() {
	req := &ActiveReq{
		UserID: "c99b9fbc-9a4e-43c0-b3e2-9c5dc55782a0",
		KeyID:  "KTdcf7HesCicFNo2ZP58jQ",
		Active: rule.ActiveDefault,
	}

	resp, err := suite.keyGen.ActiveKey(suite.ctx, req)
	assert.Nil(suite.T(), err)
	fmt.Println("=================Update===================")
	fmt.Printf("%+v\n", resp)
}

func (suite *KeyGenSuite) TestList() {
	req := &ListReq{
		Page:   1,
		Limit:  10,
		UserID: "c99b9fbc-9a4e-43c0-b3e2-9c5dc55782a0",
	}

	resp, err := suite.keyGen.ListKey(suite.ctx, req)
	assert.Nil(suite.T(), err)
	for _, v := range resp.Keys {
		fmt.Printf("%+v \n", v)
	}
}

func (suite *KeyGenSuite) TestSignature() {
	req := &SignatureReq{
		AccessKeyID: "KTdcf7HesCicFNo2ZP58jQ",
		FullBody: map[string]interface{}{
			"other": "other",
		},
	}

	resp, err := suite.keyGen.Signature(suite.ctx, req)
	assert.Nil(suite.T(), err)
	fmt.Println("================================")
	fmt.Printf("%+v\n", resp)
}

func (suite *KeyGenSuite) TestDel() {
	req := &DeleteReq{
		// ID:     "031611fe-dc7d-42ed-a95f-191d649e3ba9",
		KeyID:  "KTdcf7HesCicFNo2ZP58jQ",
		UserID: "c99b9fbc-9a4e-43c0-b3e2-9c5dc55782a0",
	}

	_, err := suite.keyGen.DeleteKey(suite.ctx, req)
	assert.Nil(suite.T(), err)
}

func (suite *KeyGenSuite) _TestKeyGen() {
	key, keySecret, encodeSecret, err := genKey()
	assert.Nil(suite.T(), err)
	fmt.Println(key)
	fmt.Println(keySecret)
	fmt.Println(encodeSecret)
}

// func (suite *KeyGenSuite) TestEncode() {
// 	encode, err := netcrypto.EasyEncodeString("abc")
// 	// 01100001 01100010 01100011
// 	assert.Nil(suite.T(), err)
// 	fmt.Println("encode: " + encode)
// }
