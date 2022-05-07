package service

import (
	"context"
	"encoding/json"
	"fmt"
	"kms/internal/apipath"
	"kms/internal/eauth"
	"kms/internal/enums"
	"kms/internal/models"
	"kms/internal/models/mysql"
	"kms/internal/models/redis"
	"kms/internal/rule"
	"kms/pkg/crypto/encode"
	"kms/pkg/dbcli"
	"kms/pkg/math"
	"kms/pkg/misc/code"
	"kms/pkg/misc/config"
	"kms/pkg/xsvc"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"gorm.io/gorm"
)

// KeyAgent external key interface
type KeyAgent interface {
	Upload(ctx context.Context, req *UploadAgencyReq) (*UploadAgencyResp, error)
	CheckAuth(ctx context.Context, req *CheckAuthReq) (*CheckAuthResp, error)
	Query(ctx context.Context, req *QueryAgencyReq) (*QueryAgencyResp, error)
	List(ctx context.Context, req *ListAgencyReq) (*ListAgencyResp, error)
	DeleteInBatch(ctx context.Context, req *DeleteInBatchReq) (*DeleteInBatchResp, error)
	Delete(ctx context.Context, req *DeleteAgencyReq) (*DeleteAgencyResp, error)
	Authorize(ctx context.Context, req *AuthReq) (*AuthResp, error)
	ListAuthType(ctx context.Context, req *ListAuthTypeReq) (*ListAuthTypeResp, error)
	GetSample(ctx context.Context, req *GetSampleReq) (*GetSampleResp, error)
	Update(ctx context.Context, req *UpdateAgencyReq) (*UpdateAgencyResp, error)
	UpdateInBatch(ctx context.Context, req *UpdateInBatchReq) (*UpdateInBatchResp, error)
	Active(ctx context.Context, req *ActiveAgencyReq) (*ActiveAgencyResp, error)
	DeleteByPrefixPath(ctx context.Context, req *DeleteByPrefixPathReq) (*DeleteByPrefixPathResp, error)
}

// CreateKeyAgent create
func CreateKeyAgent(conf *config.Config, log logger.AdaptedLogger) (KeyAgent, error) {
	db, err := dbcli.GetMysqlClient(nil, log)
	if err != nil {
		return nil, err
	}

	rc, err := dbcli.GetRedisClient(nil)
	if err != nil {
		return nil, err
	}

	return &keyAgent{
		db:         db,
		log:        log.WithName("[server] Key Agent"),
		ExtKeyRepo: mysql.CreateExtKeyRepo(),
		redis:      redis.NewRedisClient(rc),
	}, nil
}

// ********************************Upload ExtKey*********************************

type keyAgent struct {
	db         *gorm.DB
	ExtKeyRepo models.AgencyKeyRepo
	redis      models.Cache
	log        logger.AdaptedLogger
}

// UploadAgencyReq req
type UploadAgencyReq struct {
	UserID      string
	UserName    string
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description" binding:"max=255"`
	Service     string `json:"service" binding:"required"`
	Host        string `json:"host" binding:"required"`
	AuthType    string `json:"authType" binding:"required"`
	AuthContent string `json:"authContent"`
	KeyID       string `json:"keyID" binding:"max=512"`
	KeySecret   string `json:"keySecret" binding:"max=4096"`
	KeyContent  string `json:"keyContent"`
}

// UploadAgencyResp resp
type UploadAgencyResp struct {
	ID string `json:"id"`
}

// Upload upload
func (ka *keyAgent) Upload(ctx context.Context, req *UploadAgencyReq) (*UploadAgencyResp, error) {
	if err := rule.CheckCharSet(req.KeyID, req.KeySecret, req.Description); err != nil {
		return nil, err
	}

	if req.Service == "" {
		return nil, error2.NewErrorWithString(error2.ErrParams, "missing service")
	}

	if _, err := rule.ParseAuthContent(req.AuthContent, true); err != nil {
		return nil, error2.NewErrorWithString(error2.ErrParams,
			fmt.Sprintf("AuthContent error:%s", err.Error()))
	}

	entity, err := ka.ExtKeyRepo.QueryInService(ka.db, req.Service, req.KeyID)
	if err != nil {
		return nil, err
	}
	if entity.ID != "" {
		return nil, error2.New(code.ErrKeyIsExists)
	}

	secret, err := encode.SecretEncodeString(req.KeySecret, req.KeyID)
	if err != nil {
		return nil, err
	}

	curTime := time2.NowUnix()
	ak := &models.AgencyKey{
		ID:          id2.StringUUID(),
		Owner:       req.UserID,
		OwnerName:   req.UserName,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		Service:     req.Service,
		Host:        req.Host,
		AuthType:    req.AuthType,
		AuthContent: req.AuthContent,
		KeyID:       req.KeyID,
		KeySecret:   secret,
		KeyContent:  req.KeyContent,
		Active:      rule.ActiveDefault,
		Parsed:      rule.NotParsed,
		CreateAt:    curTime,
		UpdateAt:    curTime,
	}

	if err := ka.ExtKeyRepo.Create(ka.db, ak); err != nil {
		ka.log.Error("failed to create agency key: ", err.Error())
		return nil, err
	}

	if err := ka.redis.Cache(redis.CacheAgencyKey, ak.ID, ak); err != nil {
		return nil, err
	}
	ka.redis.Del(redis.CacheAgencyKeyList, req.Service, cacheSingleService(req.Service))

	return &UploadAgencyResp{
		ID: ak.ID,
	}, nil
}

// *********************************Query ExtKey*********************************

// QueryAgencyReq req
type QueryAgencyReq struct {
	ID string `json:"id" binding:"required"`
}

// QueryAgencyResp resp
type QueryAgencyResp struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`
	Host        string `json:"host"`
	AuthType    string `json:"authType"`
	AuthContent string `json:"authContent"`
	KeyID       string `json:"keyID"`
	KeyContent  string `json:"keyContent"`
	Active      int    `json:"active"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

// Query query
func (ka *keyAgent) Query(ctx context.Context, req *QueryAgencyReq) (*QueryAgencyResp, error) {
	ak, err := ka.query(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &QueryAgencyResp{
		ID:          ak.ID,
		Owner:       ak.Owner,
		OwnerName:   ak.OwnerName,
		Name:        ak.Name,
		Title:       ak.Title,
		Description: ak.Description,
		Service:     ak.Service,
		Host:        ak.Host,
		AuthType:    ak.AuthType,
		AuthContent: ak.AuthContent,
		KeyID:       ak.KeyID,
		KeyContent:  ak.KeyContent,
		Active:      ak.Active,
		CreateAt:    ak.CreateAt,
		UpdateAt:    ak.UpdateAt,
	}, nil
}

func (ka *keyAgent) query(ctx context.Context, id string) (*models.AgencyKey, error) {
	ak := &models.AgencyKey{}
	err := ka.redis.Query(redis.CacheAgencyKey, id, ak)
	if err != nil {
		sk, err := ka.ExtKeyRepo.Query(ka.db, id)
		if err != nil {
			return nil, err
		}
		if err = ka.redis.Cache(redis.CacheAgencyKey, sk.ID, sk); err != nil {
			return nil, err
		}
		ak = sk
	}
	return ak, nil
}

// **********************************List ExtKey**********************************

// ListAgencyReq req
type ListAgencyReq struct {
	Page    int    `json:"page" binding:"gt=0"`
	Limit   int    `json:"limit"`
	Service string `json:"service" binding:"required"`
	Active  int    `json:"active"` // -1 = all
	Owner   string
}

// ListAgencyResp resp
type ListAgencyResp struct {
	Keys  []*models.SimplifyAgencyKey `json:"keys"`
	Total int                         `json:"total"`
}

// List list
func (ka *keyAgent) List(ctx context.Context, req *ListAgencyReq) (*ListAgencyResp, error) {
	limitOne := req.Limit == 1
	aks := make([]*models.SimplifyAgencyKey, 0)
	if limitOne {
		err := ka.redis.Query(redis.CacheAgencyKeyList, cacheSingleService(req.Service), &aks)
		if err == nil {
			return &ListAgencyResp{
				Keys:  aks,
				Total: len(aks),
			}, nil
		}
	}

	aks, err := ka.list(ctx, req.Active, req.Service)
	if err != nil {
		return nil, err
	}

	filterKeys := filtKeys(req, aks)
	total := len(filterKeys)

	var ret []*models.SimplifyAgencyKey
	if req.Limit > 0 && len(filterKeys) > 0 {
		curPage := (req.Page - 1) * req.Limit
		nextPage := req.Page * req.Limit
		if curPage < total {
			ret = filterKeys[curPage:math.MinInt(nextPage, total)]
		}
	} else {
		ret = filterKeys
	}

	if limitOne {
		ka.redis.Cache(redis.CacheAgencyKeyList, cacheSingleService(req.Service), ret)
	}

	return &ListAgencyResp{
		Keys:  ret,
		Total: total,
	}, nil
}

func cacheSingleService(service string) string {
	return fmt.Sprintf("%s#1", service)
}

func (ka *keyAgent) list(ctx context.Context, active int, service string) ([]*models.SimplifyAgencyKey, error) {
	if service == "" {
		return nil, error2.New(error2.ErrParams)
	}

	aks := make([]*models.SimplifyAgencyKey, 0)
	err := ka.redis.Query(redis.CacheAgencyKeyList, service, &aks)
	if err != nil {
		if aks, _, err = ka.ExtKeyRepo.List(ka.db, service); err != nil {
			return nil, err
		}
		ka.redis.Cache(redis.CacheAgencyKeyList, service, aks)
	}
	return aks, nil
}

func filtKeys(req *ListAgencyReq, aks []*models.SimplifyAgencyKey) []*models.SimplifyAgencyKey {
	if req.Owner != "" || req.Active >= 0 {
		filter := aks[:0] //reuse the origin buffer
		for _, v := range aks {
			if !((req.Owner != "" && req.Owner != v.Owner) ||
				(req.Active > 0 && req.Active != v.Active)) {
				filter = append(filter, v)
			}
		}
		return filter
	}
	return aks
}

// ********************************Delete ExtKey*********************************

// DeleteAgencyReq req
type DeleteAgencyReq struct {
	ID      string `json:"id" binding:"required"`
	Service string `json:"service" binding:"required"`
	UserID  string
}

// DeleteAgencyResp resp
type DeleteAgencyResp struct {
}

// Delete delete
func (ka *keyAgent) Delete(ctx context.Context, req *DeleteAgencyReq) (*DeleteAgencyResp, error) {
	if _, err := ka.check(ctx, req.ID, OpDelete); err != nil {
		return nil, err
	}

	err := ka.ExtKeyRepo.Delete(ka.db, req.ID)
	if err != nil {
		return nil, err
	}

	ka.redis.Del(redis.CacheAgencyKey, req.ID)
	ka.redis.Del(redis.CacheAgencyKeyList, req.Service, cacheSingleService(req.Service))

	return &DeleteAgencyResp{}, nil
}

// **********************************Delete In Batch**********************************

// DeleteInBatchReq req
type DeleteInBatchReq struct {
	Namespace   string   `json:"namespace" binding:"required"`
	ServiceName []string `json:"serviceName" binding:"required"`
}

// DeleteInBatchResp resp
type DeleteInBatchResp struct{}

// DeleteInBatch DeleteInBatch
func (ka *keyAgent) DeleteInBatch(ctx context.Context, req *DeleteInBatchReq) (*DeleteInBatchResp, error) {
	servicePaths := make([]string, 0, len(req.ServiceName))
	for _, v := range req.ServiceName {
		servicePaths = append(servicePaths, apipath.Join(req.Namespace, v))
	}

	err := ka.deleteInBatch(ctx, servicePaths)
	return &DeleteInBatchResp{}, err
}

func (ka *keyAgent) deleteInBatch(ctx context.Context, servicePaths []string) error {
	ids, err := ka.ExtKeyRepo.DeleteInBatch(ka.db, servicePaths)
	if err != nil {
		return err
	}

	if err := ka.redis.Del(redis.CacheAgencyKey, ids...); err != nil {
		ka.log.Error("failed to delete agency key", err.Error())
	}
	for _, service := range servicePaths {
		ka.redis.Del(redis.CacheAgencyKeyList, service, cacheSingleService(service))
	}
	return nil
}

// DeleteByPrefixPathReq DeleteByPrefixPathReq
type DeleteByPrefixPathReq struct {
	NamespacePath string `json:"namespacePath"`
}

// DeleteByPrefixPathResp DeleteByPrefixPathResp
type DeleteByPrefixPathResp struct {
}

func (ka *keyAgent) DeleteByPrefixPath(ctx context.Context, req *DeleteByPrefixPathReq) (*DeleteByPrefixPathResp, error) {
	keys, err := ka.ExtKeyRepo.DelByPrefixPath(ka.db, req.NamespacePath)
	if err != nil {
		return nil, err
	}

	for _, v := range keys {
		ka.redis.Del(redis.CacheAgencyKey, v.ID)
		ka.redis.Del(redis.CacheAgencyKeyList, v.Service, cacheSingleService(v.Service))
	}
	return &DeleteByPrefixPathResp{}, nil
}

// ********************************** Auth **********************************

// AuthReq req
type AuthReq struct {
	ID             string                 `json:"ID" binding:"required"`
	Body           map[string]interface{} `json:"body" binding:"required"`
	APIServiceArgs string                 `json:"-"`
}

// AuthResp resp
type AuthResp struct {
	Token []*eauth.AuthResp `json:"token"`
}

// Authorize 3party auth
func (ka *keyAgent) Authorize(ctx context.Context, req *AuthReq) (*AuthResp, error) {
	ak, err := ka.check(ctx, req.ID, OpSignature)
	if err != nil {
		return nil, err
	}

	if req.APIServiceArgs != "" {
		xArgs, err := xsvc.Unmarshal(req.APIServiceArgs)
		if err != nil {
			return nil, err
		}
		ak = &models.AgencyKey{
			// Note: ignore AuthType and AuthContent in current version
			// AuthType:	xArgs.AuthType,
			// AuthContent: xArgs.AuthContent,
			Host:        xArgs.Host,
			AuthType:    ak.AuthType,
			AuthContent: ak.AuthContent,
			KeyID:       xArgs.KeyID,
			KeySecret:   xArgs.KeySecret,
			Parsed:      rule.NotParsed,
		}
	}

	at, err := eauth.GetAuthFactory().Create(ak)
	if err != nil {
		ka.log.Error("failed to create authorization method: ", err.Error())
		return nil, err
	}

	token, err := at.Invoke(req.Body)
	if err != nil {
		ka.log.Error("failed to authorize: ", err.Error())
		return nil, err
	}

	// FIXME: result of func (random, date ...) will be parsed.
	if !rule.CheckParse(ak.Parsed) && ak.ID != "" {
		b, err := json.Marshal(at.GetContent())
		if err != nil {
			return nil, err
		}
		ak.AuthContent = string(b)
		ak.Parsed = rule.Parsed
		if err := ka.ExtKeyRepo.UpdateContent(ka.db, ak.ID, ak.AuthContent, ak.Parsed); err != nil {
			return nil, err
		}
		ka.redis.Del(redis.CacheAgencyKey, ak.ID)
		ka.redis.Del(redis.CacheAgencyKeyList, ak.Service, cacheSingleService(ak.Service))
	}

	return &AuthResp{
		Token: token,
	}, nil
}

// ********************************** ListAuthTypes **********************************

// ListAuthTypeReq req
type ListAuthTypeReq struct {
}

// ListAuthTypeResp resp
type ListAuthTypeResp struct {
	Types []*ListAuthTypeNode `json:"types"`
}

// ListAuthTypeNode node
type ListAuthTypeNode struct {
	Name   string      `json:"name"`
	Sample interface{} `json:"sample"`
}

// ListAuthType list auth
func (ka *keyAgent) ListAuthType(ctx context.Context, req *ListAuthTypeReq) (*ListAuthTypeResp, error) {
	factory := eauth.GetAuthFactory()
	types := make([]*ListAuthTypeNode, 0)
	for _, v := range enums.AuthTypeSet.GetAll() {
		auth, err := factory.CreateSample(v)
		if err != nil {
			continue
		}
		types = append(types, &ListAuthTypeNode{
			Name:   v,
			Sample: auth.(eauth.Auth).GetContentVals(),
		})
	}
	return &ListAuthTypeResp{
		Types: types,
	}, nil
}

// ********************************** Get Auth Sample **********************************

// GetSampleReq GetSampleReq
type GetSampleReq struct {
	Type string `json:"type"`
}

// GetSampleResp GetSampleResp
type GetSampleResp struct {
	Sample interface{} `json:"sample"`
}

func (ka *keyAgent) GetSample(ctx context.Context, req *GetSampleReq) (*GetSampleResp, error) {
	f := eauth.GetAuthFactory()
	sample, err := f.CreateSample(req.Type)
	if err != nil {
		return nil, error2.New(code.ErrInvalidAuthType, enums.AuthTypeSet.GetAll())
	}
	return &GetSampleResp{
		Sample: sample.(eauth.Auth).GetContentVals(),
	}, nil
}

// ********************************** Update **********************************

// UpdateAgencyReq req
type UpdateAgencyReq struct {
	ID          string `json:"id" binding:"required"`
	Service     string `json:"service" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateAgencyResp resp
type UpdateAgencyResp struct{}

func (ka *keyAgent) Update(ctx context.Context, req *UpdateAgencyReq) (*UpdateAgencyResp, error) {
	err := ka.ExtKeyRepo.Update(ka.db, req.ID, req.Title, req.Description)
	if err != nil {
		return nil, err
	}
	ka.redis.Del(redis.CacheAgencyKey, req.ID)
	ka.redis.Del(redis.CacheAgencyKeyList, req.Service, cacheSingleService(req.Service))
	return &UpdateAgencyResp{}, nil
}

// ********************************** UpdateInBatch **********************************

// UpdateInBatchReq UpdateInBatchReq
type UpdateInBatchReq struct {
	Host        string `json:"host" binding:"required"`
	Service     string `json:"service" binding:"required"`
	AuthType    string `json:"authType" binding:"required"`
	AuthContent string `json:"authContent"`
}

// UpdateInBatchResp UpdateInBatchResp
type UpdateInBatchResp struct{}

// UpdateInBatch UpdateInBatch
func (ka *keyAgent) UpdateInBatch(ctx context.Context, req *UpdateInBatchReq) (*UpdateInBatchResp, error) {
	ids, err := ka.ExtKeyRepo.UpdateInBatch(ka.db, req.Host, req.Service, req.AuthType, req.AuthContent)
	if err != nil {
		return nil, err
	}

	// FIXME: remove auth result cache
	if err := ka.redis.Del(redis.CacheAgencyKeyList, req.Service, cacheSingleService(req.Service)); err != nil {
		ka.log.Error("failed to delete agencyList: ", err.Error())
	}
	if err := ka.redis.Del(redis.CacheAgencyKey, ids...); err != nil {
		ka.log.Error("failed to delete agency key: ", err.Error())
	}
	return &UpdateInBatchResp{}, nil
}

// **********************************Active agency**********************************

// ActiveAgencyReq req
type ActiveAgencyReq struct {
	ID      string `json:"id" binding:"required"`
	Active  int    `json:"active"`
	Service string `json:"service" binding:"required"`
}

// ActiveAgencyResp resp
type ActiveAgencyResp struct {
}

// Active update agency key status
func (ka *keyAgent) Active(ctx context.Context, req *ActiveAgencyReq) (*ActiveAgencyResp, error) {
	err := ka.ExtKeyRepo.Active(ka.db, req.ID, req.Active)

	ka.redis.Del(redis.CacheAgencyKey, req.ID)
	ka.redis.Del(redis.CacheAgencyKeyList, req.Service, cacheSingleService(req.Service))
	return &ActiveAgencyResp{}, err
}

func (ka *keyAgent) check(ctx context.Context, id string, op rule.Operation) (*models.AgencyKey, error) {
	key, err := ka.query(ctx, id)
	if err != nil {
		return nil, err
	}

	return key, rule.ValidateActive(key.Active, op)
}

// **********************************Check Content**********************************

// CheckAuthReq CheckAuthReq
type CheckAuthReq struct {
	AuthType    string `json:"authType"`
	ServicePath string `json:"servicePath"`
	AuthContent string `json:"authContent"`
}

// CheckAuthResp CheckAuthResp
type CheckAuthResp struct {
}

func (ka *keyAgent) CheckAuth(ctx context.Context, req *CheckAuthReq) (*CheckAuthResp, error) {
	if err := validateAuthType(req.AuthType); err != nil {
		return nil, err
	}

	switch enums.Enum(req.AuthType) {
	case enums.AuthNone:
		key, err := ka.list(ctx, rule.ActiveAny, req.ServicePath)
		if err != nil {
			return nil, err
		}
		if len(key) > 0 {
			return nil, error2.New(code.ErrServiceWithKeys)
		}
	default:
		if _, err := rule.ParseAuthContent(req.AuthContent, true); err != nil {
			return nil, error2.New(code.ErrAuthContent, err.Error())
		}
	}

	return &CheckAuthResp{}, nil
}

func validateAuthType(at string) error {
	if !enums.AuthTypeSet.Verify(at) {
		return error2.New(code.ErrInvalidAuthType, enums.AuthTypeSet.GetAll())
	}
	return nil
}
