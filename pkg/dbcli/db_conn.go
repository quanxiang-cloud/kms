package dbcli

import (
	"fmt"
	"kms/pkg/misc/config"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"gorm.io/gorm"
)

var (
	mysqlConn *gorm.DB
	redisConn *redis.ClusterClient
)

// InitDB initDB
func InitDB(conf *config.Config, log logger.AdaptedLogger) error {
	var err error
	mysqlConn, err = GetMysqlClient(conf, log)
	if err != nil {
		return err
	}
	redisConn, err = GetRedisClient(conf)
	if err != nil {
		return err
	}
	return nil
}

// GetMysqlClient get mysql client
func GetMysqlClient(conf *config.Config, log logger.AdaptedLogger) (*gorm.DB, error) {
	if conf != nil {
		db, err := mysql2.New(conf.Mysql, log)
		return db, err
	}
	if mysqlConn != nil {
		return mysqlConn, nil
	}
	return nil, fmt.Errorf("not init client")
}

// GetRedisClient get redis client
func GetRedisClient(conf *config.Config) (*redis.ClusterClient, error) {
	if conf != nil {
		db, err := redis2.NewClient(conf.Redis)
		return db, err
	}
	if redisConn != nil {
		return redisConn, nil
	}
	return nil, fmt.Errorf("not init client")
}
