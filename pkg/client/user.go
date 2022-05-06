package client

import (
	"context"
	"kms/pkg/misc/config"
	"net/http"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	usersInfoURI = "http://org/api/v1/org/o/user/ids"
)

// User user
type User interface {
	GetUsersInfo(ctx context.Context, req *GetUsersInfoReq) (*GetUsersInfoResp, error)
}

type user struct {
	Client http.Client
}

// NewUser new
func NewUser(c *config.Config) User {
	return &user{
		Client: client.New(client.Config{
			Timeout:      c.InternalNet.Timeout,
			MaxIdleConns: c.InternalNet.MaxIdleConns,
		}),
	}
}

// GetUsersInfoReq GetUsersInfoReq
type GetUsersInfoReq struct {
	IDs []string `json:"ids"`
}

// GetUsersInfoResp GetUsersInfoResp
type GetUsersInfoResp struct {
	Users []*UserInfo `json:"users"`
}

// UserInfo UserInfo
type UserInfo struct {
	ID        string        `json:"id,omitempty" `
	Name      string        `json:"name,omitempty" `
	Phone     string        `json:"phone,omitempty" `
	Email     string        `json:"email,omitempty" `
	SelfEmail string        `json:"selfEmail,omitempty" `
	UseStatus int           `json:"useStatus,omitempty" ` //状态：1正常，-2禁用，-3离职，-1删除，2激活==1 （与账号库相同）
	TenantID  string        `json:"tenantID,omitempty" `  //租户id
	Position  string        `json:"position,omitempty" `  //职位
	Avatar    string        `json:"avatar,omitempty" `    //头像
	JobNumber string        `json:"jobNumber,omitempty" ` //工号
	Status    int           `json:"status"`               //第一位：密码是否需要重置
	Dep       [][]*DeptInfo `json:"deps,omitempty"`       //用户所在部门
	Leader    [][]*UserInfo `json:"leaders,omitempty"`    //用户所在部门
}

// DeptInfo DeptInfo
type DeptInfo struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	LeaderID  string `json:"leaderID"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid"`
	SuperPID  string `json:"superID,omitempty"`
	Grade     int    `json:"grade,omitempty"`
	Attr      int    `json:"attr"`
}

// GetUsersInfo GetUsersInfo
func (u *user) GetUsersInfo(ctx context.Context, req *GetUsersInfoReq) (*GetUsersInfoResp, error) {
	resp := &GetUsersInfoResp{}
	err := client.POST(ctx, &u.Client, usersInfoURI, req, resp)
	return resp, err
}
