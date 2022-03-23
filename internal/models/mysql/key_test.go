package mysql

import (
	"context"
	"fmt"
	"kms/internal/models"
	"kms/pkg/misc/config"
	"strconv"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type keySuite struct {
	suite.Suite
	db            *gorm.DB
	ctx           context.Context
	keyRepo       models.KeyRepo
	keyConfigRepo models.KeyConfigRepo
}

func _TestKey(t *testing.T) {
	suite.Run(t, new(keySuite))
}

func (suite *keySuite) SetupSuite() {
	config, err := config.NewConfig("./testdata/config.yml")
	assert.Nil(suite.T(), err)

	log := logger.New(&config.Logger)

	db, _ := mysql2.New(config.Mysql, log)
	suite.db = db
	suite.keyRepo = NewKeyRepo()
	suite.keyConfigRepo = NewKeyConfigRepo()
}

func (suite *keySuite) TestTime() {
	system, err := suite.keyConfigRepo.Query(suite.db, "system")
	assert.Nil(suite.T(), err)
	_ = system
	fmt.Println(system.ConfigContent)
	fmt.Println(system.CreateAt)
}

func (suite *keySuite) _TestCreate() {
	key := &models.Key{
		KeyID:       "keyID2",
		Name:        "name",
		Description: "description",
		KeySecret:   "secret",
		Active:      1,
		Owner:       "userID",
	}
	err := suite.keyRepo.Create(suite.db, key)
	assert.Nil(suite.T(), err)
}

func (suite *keySuite) _TestQuery() {
	key, err := suite.keyRepo.Query(suite.db, "keyID2")
	assert.Nil(suite.T(), err)
	fmt.Printf("%+v\n", key)
}

func (suite *keySuite) _TestList() {
	keys, total, err := suite.keyRepo.List(suite.db, "")
	assert.Nil(suite.T(), err)
	for _, v := range keys {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println("total: " + strconv.Itoa(int(total)))
}

func (suite *keySuite) _TestDelete() {
	err := suite.keyRepo.Delete(suite.db, "keyID")
	assert.Nil(suite.T(), err)
}
