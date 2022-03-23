package redis

import "time"

const (
	keyPrefix      = "kms"
	defaultTimeout = time.Minute * 5
)

// CacheType CacheType
type CacheType = cacheType
type cacheType int

// CacheType
const (
	CacheInnerKey cacheType = iota
	CacheInnerKeyList

	CacheAgencyKey
	CacheAgencyKeyList

	CacheOAuth2Config
	CacheOAuth2Token

	CacheCookies

	CacheKeyConfig
)

func typeof(ct cacheType) string {
	switch ct {
	case CacheInnerKey:
		return "k"
	case CacheInnerKeyList:
		return "kl"
	case CacheAgencyKey:
		return "ak"
	case CacheAgencyKeyList:
		return "akl"
	case CacheOAuth2Config:
		return "oa"
	case CacheOAuth2Token:
		return "oat"
	case CacheCookies:
		return "ck"
	case CacheKeyConfig:
		return "kc"
	default:
		return "?"
	}
}
