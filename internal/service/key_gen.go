package service

import (
	"context"
	"kms/internal/adaptor"
	"kms/internal/enums"
	"kms/internal/models"
	"kms/internal/models/mysql"
	"kms/internal/models/redis"
	"kms/internal/rule"
	client "kms/pkg/client"
	"kms/pkg/crypto/decode"
	"kms/pkg/crypto/encode"
	"kms/pkg/dbcli"
	"kms/pkg/hash"
	"kms/pkg/math"
	"kms/pkg/misc/code"
	"kms/pkg/misc/config"
	"kms/pkg/signature"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"gorm.io/gorm"
)

// KeyGenerator kms
type KeyGenerator interface {
	CreateKey(ctx context.Context, req *CreateReq) (*CreateResp, error)
	DeleteKey(ctx context.Context, req *DeleteReq) (*DeleteResp, error)
	ListKey(ctx context.Context, req *ListReq) (*ListResp, error)
	QueryKey(ctx context.Context, req *QueryReq) (*QueryResp, error)
	ActiveKey(ctx context.Context, req *ActiveReq) (*ActiveResp, error)
	UpdateKey(ctx context.Context, req *UpdateReq) (*UpdateResp, error)
	Signature(ctx context.Context, req *SignatureReq) (*SignatureResp, error)
}

// CreateKeyGenerator create keyGenerator
func CreateKeyGenerator(conf *config.Config, log logger.AdaptedLogger) (KeyGenerator, error) {
	db, err := dbcli.GetMysqlClient(nil, log)
	if err != nil {
		return nil, err
	}

	rc, err := dbcli.GetRedisClient(nil)
	if err != nil {
		return nil, err
	}

	k := &keyGenerator{
		db:      db,
		log:     log.WithName("[server] key"),
		org:     client.NewUser(conf),
		redis:   redis.NewRedisClient(rc),
		keyRepo: mysql.NewKeyRepo(),
	}

	return k, nil
}

// **********************************Create Key**********************************

type keyGenerator struct {
	db      *gorm.DB
	log     logger.AdaptedLogger
	org     client.User
	redis   models.Cache
	keyRepo models.KeyRepo
	// keyConfigRepo models.KeyConfigRepo
}

// CreateReq req
type CreateReq struct {
	UserID      string
	UserName    string
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
}

// CreateResp resp
type CreateResp struct {
	ID        string `json:"id"`
	KeyID     string `json:"keyID"`
	KeySecret string `json:"keySecret"`
}

// CreateKey create key
func (k *keyGenerator) CreateKey(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	if err := rule.CheckCharSet(req.Title, req.Description); err != nil {
		return nil, err
	}

	err := checkNum(ctx, k, req.UserID)
	if err != nil {
		return nil, err
	}

	keyID, keySecret, encodeSecret, err := genKey()
	if err != nil {
		k.log.Error("failed to encode key: ", err.Error())
		return nil, err
	}

	curTime := time2.NowUnix()
	key := &models.Key{
		ID:          id2.StringUUID(),
		Owner:       req.UserID,
		OwnerName:   req.UserName,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		KeyID:       keyID,
		KeySecret:   encodeSecret,
		Active:      rule.ActiveDefault,
		Assignee:    req.Assignee,
		CreateAt:    curTime,
		UpdateAt:    curTime,
	}

	err = k.keyRepo.Create(k.db, key)
	if err != nil {
		return nil, err
	}

	// cache in memory
	k.redis.Cache(redis.CacheInnerKey, key.ID, key)

	// use elimination strategy
	k.redis.Del(redis.CacheInnerKeyList, req.UserID)

	return &CreateResp{
		ID:        key.ID,
		KeyID:     key.KeyID,
		KeySecret: keySecret,
	}, nil
}

func checkNum(ctx context.Context, k *keyGenerator, userID string) error {
	exist, err := k.keyRepo.Count(k.db, userID)
	if err != nil {
		return err
	}

	limit, err := getLimit(ctx, userID)
	if err != nil {
		return err
	}

	if exist >= limit {
		return error2.New(code.ErrMaximumHold)
	}
	return nil
}

func getLimit(ctx context.Context, userID string) (int, error) {
	kcOper, err := adaptor.GetKeyConfigOper()
	if err != nil {
		return 0, err
	}

	limit, err := kcOper.Query(ctx, &adaptor.QueryKConfigReq{
		Owner: userID,
	})
	if err != nil {
		return -1, err
	}

	if limit != nil && limit.ID != "" {
		keyNumCfg, ok := limit.ConfigContent[enums.ConfigKeyNum.Val()]
		if ok {
			expiryCfg := limit.ConfigContent[enums.ConfigKeyExpiry.Val()]
			if _, err := rule.CheckKeyExpiry(expiryCfg); err == nil {
				return rule.CheckKeyNum(keyNumCfg)
			}
		}
	}
	return models.DefaultNum, nil
}

// gen key and secret(encoded by key)
func genKey() (keyID, keySecret, encodeSecret string, err error) {
	uuid := id2.StringUUID()
	keyID = hash.Md5Hash(1024, uuid)
	uuid = id2.StringUUID()
	keySecret = hash.Sha256Hash(1024, uuid)
	encodeSecret, err = encode.SecretEncodeString(keySecret, keyID)
	return
}

// **********************************Delete Key**********************************

// DeleteReq req
type DeleteReq struct {
	KeyID  string `json:"keyID" binding:"required"`
	UserID string
}

// DeleteResp resp
type DeleteResp struct {
}

// DeleteKey del key
func (k *keyGenerator) DeleteKey(ctx context.Context, req *DeleteReq) (*DeleteResp, error) {
	if _, err := k.check(ctx, req.KeyID, OpDelete); err != nil {
		return nil, err
	}

	affected := k.keyRepo.Delete(k.db, req.KeyID)
	if affected < 1 {
		return nil, error2.New(code.ErrKeyDelFaild)
	}

	k.redis.Del(redis.CacheInnerKey, req.KeyID)
	k.redis.Del(redis.CacheInnerKeyList, req.UserID)

	return &DeleteResp{}, nil
}

// ********************************** List Key **********************************

// ListReq req
type ListReq struct {
	Page   int `json:"page" binding:"gt=0"`
	Limit  int `json:"limit"`
	UserID string
}

// ListResp resp
type ListResp struct {
	Keys  []*models.SimplifyKey `json:"keys"`
	Total int                   `json:"total"`
}

// ListKey list
func (k *keyGenerator) ListKey(ctx context.Context, req *ListReq) (*ListResp, error) {
	var keys []*models.SimplifyKey
	var total int
	var err error

	err = k.redis.Query(redis.CacheInnerKeyList, req.UserID, &keys)
	total = len(keys)
	if err != nil {
		keys, total, err = k.keyRepo.List(k.db, req.UserID)
		if err != nil {
			return nil, err
		}
		k.redis.Cache(redis.CacheInnerKeyList, req.UserID, keys)
	}

	ret := make([]*models.SimplifyKey, 0)
	if req.Limit > 0 {
		curPage := (req.Page - 1) * req.Limit
		nextPage := req.Page * req.Limit
		if curPage < total {
			ret = keys[curPage:math.MinInt(nextPage, total)]
		}
	} else {
		ret = keys
	}

	return &ListResp{
		Keys:  ret,
		Total: total,
	}, nil
}

// ********************************* Query Key **********************************

// QueryReq queryReq
type QueryReq struct {
	KeyID string `json:"keyID" binding:"required"`
}

// QueryResp queryResp
type QueryResp struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	KeyID       string `json:"keyID"`
	Active      int    `json:"active"`
	Assignee    string `json:"assignee"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

// Query query key by user-id
func (k *keyGenerator) QueryKey(ctx context.Context, req *QueryReq) (*QueryResp, error) {
	key, err := k.query(ctx, req.KeyID)
	if err != nil {
		return nil, err
	}

	return &QueryResp{
		ID:          key.ID,
		Owner:       key.Owner,
		OwnerName:   key.OwnerName,
		Name:        key.Name,
		Title:       key.Title,
		Description: key.Description,
		KeyID:       key.KeyID,
		Active:      key.Active,
		Assignee:    key.Assignee,
		CreateAt:    key.CreateAt,
		UpdateAt:    key.UpdateAt,
	}, nil
}

func (k *keyGenerator) query(ctx context.Context, keyID string) (*models.Key, error) {
	key := &models.Key{}
	err := k.redis.Query(redis.CacheInnerKey, keyID, key)
	if err != nil {
		key, err = k.keyRepo.Query(k.db, keyID)
		if err != nil {
			return nil, err
		}
		k.redis.Cache(redis.CacheInnerKey, key.ID, key)
	}
	return key, err
}

// **********************************Active Key**********************************

// ActiveReq req
type ActiveReq struct {
	KeyID  string `json:"keyID" binding:"required"`
	Active int    `json:"active"`
	UserID string
}

// ActiveResp resp
type ActiveResp struct {
}

// UpdateKey update key status
func (k *keyGenerator) ActiveKey(ctx context.Context, req *ActiveReq) (*ActiveResp, error) {
	// Update db
	uk := &models.Key{
		KeyID:  req.KeyID,
		Active: req.Active,
	}
	if err := k.keyRepo.Active(k.db, uk); err != nil {
		return nil, err
	}

	// remove cache
	k.redis.Del(redis.CacheInnerKey, req.KeyID)
	k.redis.Del(redis.CacheInnerKeyList, req.UserID)

	return &ActiveResp{}, nil
}

// **********************************Update Key**********************************

// UpdateReq UpdateReq
type UpdateReq struct {
	UserID      string
	KeyID       string `json:"keyID" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
}

// UpdateResp UpdateResp
type UpdateResp struct {
}

func (k *keyGenerator) UpdateKey(ctx context.Context, req *UpdateReq) (*UpdateResp, error) {
	uk := &models.Key{
		KeyID:       req.KeyID,
		Title:       req.Title,
		Description: req.Description,
		Assignee:    req.Assignee,
	}
	if err := k.keyRepo.Update(k.db, uk); err != nil {
		return nil, err
	}

	// remove cache
	k.redis.Del(redis.CacheInnerKey, req.KeyID)
	k.redis.Del(redis.CacheInnerKeyList, req.UserID)

	return &UpdateResp{}, nil
}

// ********************************** Signature *********************************

// SignatureReq req
type SignatureReq struct {
	AccessKeyID string
	FullBody    map[string]interface{}
}

// SignatureResp signatureResp
type SignatureResp struct {
	Sign     string   `json:"sign"`
	UserInfo UserInfo `json:"userInfo"`
}

// UserInfo userinfo
type UserInfo struct {
	UserID        string   `json:"userID"`
	UserName      string   `json:"userName"`
	DepartmentIDs []string `json:"departmentIDs"`
}

// Signature sign keyid and gen secret
func (k *keyGenerator) Signature(ctx context.Context, req *SignatureReq) (*SignatureResp, error) {
	key, err := k.check(ctx, req.AccessKeyID, OpSignature)
	if err != nil {
		return nil, err
	}

	if key.ID == "" {
		return nil, error2.New(code.ErrInvalidKey)
	}

	// decode secret
	secret, err := decode.SecretDecodeString(key.KeySecret, key.KeyID)
	if err != nil {
		return nil, err
	}

	// signature
	sign, err := signature.Signature(req.FullBody, secret)
	if err != nil {
		return nil, err
	}

	infos, err := getUserInfo(ctx, k, key.Owner)
	if err != nil {
		return nil, err
	}
	if len(infos) <= 0 {
		return nil, error2.NewErrorWithString(error2.Internal, "not found users info")
	}
	deptIDs := make([]string, 0, 8)
	for _, dep := range infos[0].Dep {
		for _, v := range dep {
			deptIDs = append(deptIDs, v.ID)
		}
	}

	return &SignatureResp{
		Sign: sign,
		UserInfo: UserInfo{
			UserID:        infos[0].ID,
			UserName:      infos[0].Name,
			DepartmentIDs: deptIDs,
		},
	}, nil
}

// remote call to get user info
func getUserInfo(ctx context.Context, k *keyGenerator, userID string) ([]*client.UserInfo, error) {
	// remote call
	req := &client.GetUsersInfoReq{
		IDs: []string{userID},
	}

	resp, err := k.org.GetUsersInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Users) < 1 {
		return nil, error2.NewErrorWithString(error2.ErrParams, "Invalid Access-Key-ID")
	}
	return resp.Users, nil
}

func (k *keyGenerator) check(ctx context.Context, keyID string, op rule.Operation) (*models.Key, error) {
	key, err := k.query(ctx, keyID)
	if err != nil {
		return nil, err
	}
	return key, rule.ValidateActive(key.Active, op)
}
