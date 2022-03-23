package service

import (
	"context"
	"fmt"
	"kms/internal/adaptor"
	"kms/internal/enums"
	"kms/internal/models"
	"kms/internal/models/mysql"
	"kms/internal/models/redis"
	"kms/internal/rule"
	"kms/pkg/dbcli"
	"kms/pkg/misc/config"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/time"
	"gorm.io/gorm"
)

const (
	// TimeFormat time format
	TimeFormat = "2006.01.02 15:04.05"
	// SystemConfig system config
	SystemConfig = "system"
)

// KeyConfigServe key config center
type KeyConfigServe interface {
	Update(ctx context.Context, req *UpdateKConfigReq) (*UpdateKConfigResp, error)
	Query(ctx context.Context, req *QueryKConfigReq) (*QueryKConfigResp, error)
	Delete(ctx context.Context, req *DeleteKConfigReq) (*DeleteKConfigResp, error)
}

// NewKeyConfigServe new
func NewKeyConfigServe(c *config.Config, log logger.AdaptedLogger) (KeyConfigServe, error) {
	db, err := dbcli.GetMysqlClient(nil, log)
	if err != nil {
		return nil, err
	}

	rc, err := dbcli.GetRedisClient(nil)
	if err != nil {
		return nil, err
	}

	kc := &keyConfig{
		db:     db,
		log:    log.WithName("[server] key config"),
		kcRepo: mysql.NewKeyConfigRepo(),
		redis:  redis.NewRedisClient(rc),
	}

	adaptor.SetKeyConfigOper(kc)

	err = initSystemConfig(kc)
	if err != nil {
		kc.log.Error("failed to init system config: ", err.Error())
		return nil, err
	}
	return kc, nil
}

type keyConfig struct {
	db     *gorm.DB
	log    logger.AdaptedLogger
	kcRepo models.KeyConfigRepo
	redis  models.Cache
}

func initSystemConfig(kc *keyConfig) error {
	kc.log.Info("init system config")
	if models.DefaultNum < 0 {
		system, err := kc.Query(context.Background(), &adaptor.QueryKConfigReq{
			Owner: SystemConfig,
		})
		if err != nil {
			return err
		}
		keyNumConfig, ok := system.ConfigContent[enums.ConfigKeyNum.Val()]
		if !ok {
			kc.log.Error("failed to init system config: cause the default key num not exists")
			return fmt.Errorf("failed to init system config")
		}
		models.DefaultNum, err = rule.CheckKeyNum(keyNumConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// ********************************** Update **********************************

// UpdateKConfigReq UpdateKConfigReq
type UpdateKConfigReq struct {
	Owner         string            `json:"owner"`
	OwnerName     string            `json:"ownerName"`
	ConfigContent map[string]string `json:"configContent"`
}

// UpdateKConfigResp UpdateKConfigResp
type UpdateKConfigResp struct {
}

// Update Update
func (kc *keyConfig) Update(ctx context.Context, req *UpdateKConfigReq) (*UpdateKConfigResp, error) {
	if err := rule.CheckKeyConfig(req.ConfigContent); err != nil {
		return nil, err
	}

	keyConfig, err := kc.kcRepo.Query(kc.db, req.Owner)
	if err != nil {
		return nil, err
	}

	keyConfig.ConfigContent = req.ConfigContent

	keyConfig.UpdateAt = time.NowUnix()
	if keyConfig.ID == "" {
		keyConfig.ID = id2.StringUUID()
		keyConfig.Owner = req.Owner
		keyConfig.OwnerName = req.OwnerName
		keyConfig.CreateAt = keyConfig.UpdateAt
		kc.kcRepo.Create(kc.db, keyConfig)
	} else {
		kc.kcRepo.Update(kc.db, keyConfig)
	}

	kc.redis.Del(redis.CacheKeyConfig, req.Owner)

	return &UpdateKConfigResp{}, nil
}

// ********************************** Query **********************************

// QueryKConfigReq QueryKConfigReq
type QueryKConfigReq = adaptor.QueryKConfigReq

// QueryKConfigResp QueryKConfigResp
type QueryKConfigResp = adaptor.QueryKConfigResp

// Query Query
func (kc *keyConfig) Query(ctx context.Context, req *QueryKConfigReq) (*QueryKConfigResp, error) {
	keyConfig := &models.KeyConfig{}
	err := kc.redis.Query(redis.CacheKeyConfig, req.Owner, keyConfig)
	if err != nil {
		keyConfig, err = kc.kcRepo.Query(kc.db, req.Owner)
		if err != nil {
			return nil, err
		}
		err = kc.redis.Cache(redis.CacheKeyConfig, req.Owner, keyConfig)
		if err != nil {
			kc.log.Errorf("failed to cache keyConfig: %s", err.Error())
		}
	}

	return &QueryKConfigResp{
		ID:            keyConfig.ID,
		Owner:         keyConfig.Owner,
		OwnerName:     keyConfig.OwnerName,
		ConfigContent: keyConfig.ConfigContent,
		CreateAt:      keyConfig.CreateAt,
	}, nil
}

// ********************************** Delete **********************************

// DeleteKConfigReq DeleteKConfigReq
type DeleteKConfigReq struct {
	Owner string `json:"owner"`
}

// DeleteKConfigResp DeleteKConfigResp
type DeleteKConfigResp struct {
}

// Delete delete
func (kc *keyConfig) Delete(ctx context.Context, req *DeleteKConfigReq) (*DeleteKConfigResp, error) {
	err := kc.kcRepo.Delete(kc.db, req.Owner)
	if err != nil {
		return nil, err
	}

	kc.redis.Del(redis.CacheKeyConfig, req.Owner)
	return &DeleteKConfigResp{}, err
}
