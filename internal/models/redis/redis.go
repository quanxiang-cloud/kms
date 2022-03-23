package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/go-redis/redis/v8"
)

func getCacheKey(ct cacheType, key string) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefix, typeof(ct), key)
}

// NewRedisClient New
func NewRedisClient(client *redis.ClusterClient) *Client {
	return &Client{
		client: client,
	}
}

// Client Client
type Client struct {
	client *redis.ClusterClient
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Cache Cache
func (c *Client) Cache(ct cacheType, key string, cached interface{}, duration ...time.Duration) error {
	key = getCacheKey(ct, key)
	b, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	dur := defaultTimeout
	if len(duration) != 0 {
		dur = duration[0]
	}
	return c.client.Set(context.Background(), key, b, dur).Err()
}

// Del del
func (c *Client) Del(ct cacheType, key ...string) error {
	for _, v := range key {
		key := getCacheKey(ct, v)
		err := c.client.Del(context.Background(), key).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// Query query
func (c *Client) Query(ct cacheType, key string, entity interface{}, duration ...time.Duration) error {
	key = getCacheKey(ct, key)
	content, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	dur := defaultTimeout
	if len(duration) != 0 {
		dur = duration[0]
	}
	c.client.Expire(context.Background(), key, dur)
	json.Unmarshal(unsafeStringBytes(content), entity)
	return nil
}
