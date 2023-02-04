package redis

import (
	"context"
	"estimation-service/pkg"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CountUsersInSegment_whenSegmentAlreadyCached(t *testing.T) {
	asserts := assert.New(t)
	c := context.Background()
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("sport").SetVal("20")
	redisClient := ClientImpl{
		rdb:      db,
		duration: time.Hour,
	}

	res, err := redisClient.CountUsersInSegment(c, "sport")
	asserts.Nil(err)
	asserts.IsType(&pkg.UserCountBySegmentResponse{}, res)

}

func Test_CountUsersInSegment_whenSegmentNotInCache(t *testing.T) {
	asserts := assert.New(t)
	c := context.Background()
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("sport").RedisNil()
	redisClient := ClientImpl{
		rdb:      db,
		duration: time.Hour,
	}
	_, err := redisClient.CountUsersInSegment(c, "sport")
	asserts.NotNil(err)
}
