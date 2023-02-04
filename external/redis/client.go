package redis

import (
	"context"
	"estimation-service/pkg"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type ClientImpl struct {
	rdb      *redis.Client
	duration time.Duration
}

func (c *ClientImpl) CountUsersInSegment(context context.Context, title string) (*pkg.UserCountBySegmentResponse, error) {
	count, rErr := c.rdb.Get(context, title).Result()
	if rErr != nil {
		return nil, rErr
	}
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, err
	}
	return &pkg.UserCountBySegmentResponse{Segmentation: title, Count: intCount}, nil
}

func (c *ClientImpl) StoreUserCountInSegment(context context.Context, data pkg.UserCountBySegmentResponse) {
	c.rdb.Set(context, data.Segmentation, data.Count, c.duration)
}

func (c *ClientImpl) ClearUserCacheInSegment(context context.Context, title string) {
	c.rdb.Del(context, title)
}
