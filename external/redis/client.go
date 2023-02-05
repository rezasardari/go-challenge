package redis

import (
	"context"
	"estimation-service/pkg"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type ClientImpl struct {
	Rdb      *redis.Client
	Duration time.Duration
}

func (c *ClientImpl) CountUsersInSegment(context context.Context, title string) (*pkg.UserCountBySegmentResponse, error) {
	count, rErr := c.Rdb.Get(context, title).Result()
	if rErr != nil {
		return nil, rErr
	}
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, err
	}
	return &pkg.UserCountBySegmentResponse{Segment: title, Count: intCount}, nil
}

func (c *ClientImpl) StoreUserCountInSegment(context context.Context, data pkg.UserCountBySegmentResponse) {
	c.Rdb.Set(context, data.Segment, data.Count, c.Duration)
}

func (c *ClientImpl) ClearUserCacheInSegment(context context.Context, title string) {
	c.Rdb.Del(context, title)
}

func (c *ClientImpl) StoreUserCountsInAllSegment(ctx context.Context, data []pkg.UserCountBySegmentResponse) error {
	pipe := c.Rdb.Pipeline()
	for _, item := range data {
		// Queue commands in the pipeline
		pipe.Set(ctx, item.Segment, item.Count, c.Duration)
	}
	// Execute all commands in the pipeline
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
