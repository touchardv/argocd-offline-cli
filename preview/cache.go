package preview

import (
	"context"
	"time"

	"github.com/argoproj/argo-cd/v2/reposerver/cache"
	cacheutil "github.com/argoproj/argo-cd/v2/util/cache"
)

type NoopCacheClient struct{}

func NewNoopCache() *cache.Cache {
	c := cacheutil.Cache{}
	c.SetClient(&NoopCacheClient{})
	noTimeout := 0 * time.Second
	return cache.NewCache(&c, noTimeout, noTimeout, noTimeout)
}

func (c *NoopCacheClient) Set(item *cacheutil.Item) error {
	return nil
}

func (c *NoopCacheClient) Rename(oldKey string, newKey string, expiration time.Duration) error {
	return nil
}

func (c *NoopCacheClient) Get(key string, obj interface{}) error {
	return nil
}

func (c *NoopCacheClient) Delete(key string) error {
	return nil
}

func (c *NoopCacheClient) OnUpdated(ctx context.Context, key string, callback func() error) error {
	return nil
}

func (c *NoopCacheClient) NotifyUpdated(key string) error {
	return nil
}
