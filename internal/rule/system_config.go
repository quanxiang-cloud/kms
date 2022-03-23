package rule

import (
	"fmt"
	"kms/internal/enums"
	"strconv"

	time2 "github.com/quanxiang-cloud/cabin/time"
)

// CheckKeyConfig checkSystemConfig
func CheckKeyConfig(cfg map[string]string) (err error) {
	for k, v := range cfg {
		if !enums.SystemConfigSet.Verify(k) {
			return fmt.Errorf("not support config type(%s), only support: %v", k, enums.SystemConfigSet.GetAll())
		}
		switch k {
		case enums.ConfigKeyNum.Val():
			_, err = CheckKeyNum(v)
		case enums.ConfigKeyExpiry.Val():
			_, err = CheckKeyExpiry(v)
		}
	}
	return
}

// CheckKeyNum CheckKeyNum
func CheckKeyNum(cfg string) (int, error) {
	return strconv.Atoi(cfg)
}

// CheckKeyExpiry CheckKeyExpiry
func CheckKeyExpiry(cfg string) (int64, error) {
	t, err := strconv.ParseInt(cfg, 10, 64)
	if err != nil {
		return 0, err
	}
	if t < time2.NowUnix() {
		return 0, fmt.Errorf("invaild time, its expired")
	}
	return t, nil
}
