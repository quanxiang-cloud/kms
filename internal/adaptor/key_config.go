package adaptor

import (
	"context"
	"fmt"
)

// KeyConfigOper oper
type KeyConfigOper interface {
	Query(ctx context.Context, req *QueryKConfigReq) (*QueryKConfigResp, error)
}

const (
	keyConfigOper = "keyConfigOper"
)

// QueryKConfigReq QueryKConfigReq
type QueryKConfigReq struct {
	Owner string `json:"owner"`
}

// QueryKConfigResp QueryKConfigResp
type QueryKConfigResp struct {
	ID            string            `json:"id"`
	Owner         string            `json:"owner"`
	OwnerName     string            `json:"ownerName"`
	ConfigContent map[string]string `json:"configContent"`
	CreateAt      int64             `json:"createAt"`
}

// SetKeyConfigOper set adaptor instance
func SetKeyConfigOper(k KeyConfigOper) {
	setInstance(keyConfigOper, k)
}

// GetKeyConfigOper get instance
func GetKeyConfigOper() (KeyConfigOper, error) {
	ins, err := getInstance(keyConfigOper)
	if err != nil {
		return nil, err
	}
	if k, ok := ins.(KeyConfigOper); ok {
		return k, nil
	}
	return nil, fmt.Errorf("instance incorrect")
}
