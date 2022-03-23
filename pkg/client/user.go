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
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	JobNumber string      `json:"jobNumber"`
	Status    int         `json:"status"`
	DepIDs    []*DeptInfo `json:"depIDs"`
	Position  string      `json:"position"`
	TenantID  string      `json:"tenantID"`
}

// DeptInfo DeptInfo
type DeptInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LeaderID string `json:"leaderID"`
	PID      string `json:"pid"`
	SuperID  string `json:"superID"`
	Grade    int    `json:"grade"`
	Attr     int    `json:"attr"`
}

// GetUsersInfo GetUsersInfo
func (u *user) GetUsersInfo(ctx context.Context, req *GetUsersInfoReq) (*GetUsersInfoResp, error) {
	resp := &GetUsersInfoResp{}
	err := client.POST(ctx, &u.Client, usersInfoURI, req, resp)
	return resp, err
}
