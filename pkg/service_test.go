package pkg

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) CountUsersInSegment(ctx context.Context, req UserCountBySegmentRequest) (*UserCountBySegmentResponse, error) {
	args := r.Called(ctx, req)
	return args.Get(0).(*UserCountBySegmentResponse), args.Error(1)
}

func (r *RepositoryMock) StoreUserInSegment(ctx context.Context, req StoreUserSegmentRequest) error {
	args := r.Called(ctx, req)
	return args.Error(0)
}

func (r *RepositoryMock) ArchiveExpiredData(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}

func (r *RepositoryMock) CountUsersInAllSegments(ctx context.Context) ([]UserCountBySegmentResponse, error) {
	args := r.Called(ctx)
	return args.Get(0).([]UserCountBySegmentResponse), args.Error(1)
}

type RedisMock struct {
	mock.Mock
	wg        sync.WaitGroup
	CallCount int
}

func (c *RedisMock) currentCount() int {
	c.wg.Wait() // wait for all the call to happen. This will block until wg.Done() is called.
	return c.CallCount
}

func (c *RedisMock) CountUsersInSegment(context context.Context, title string) (*UserCountBySegmentResponse, error) {
	args := c.Called(context, title)
	return args.Get(0).(*UserCountBySegmentResponse), args.Error(1)
}
func (c *RedisMock) StoreUserCountInSegment(context context.Context, data UserCountBySegmentResponse) {
	c.Called(context, data)
	return
}
func (c *RedisMock) ClearUserCacheInSegment(ctx context.Context, title string) {
	c.wg.Done()
	c.CallCount += 1
	c.Called(ctx, title)
	return
}
func (c *RedisMock) StoreUserCountsInAllSegment(ctx context.Context, data []UserCountBySegmentResponse) error {
	args := c.Called(ctx, data)
	return args.Error(0)
}

func Test_GetUserCountBySegmentation_WhenCacheNotNull(t *testing.T) {
	asserts := assert.New(t)

	redisMock := new(RedisMock)
	repositoryMock := new(RepositoryMock)

	sut := ServiceImpl{
		Repository: repositoryMock,
		Redis:      redisMock,
	}

	redisMock.On("CountUsersInSegment", mock.Anything, mock.Anything).Return(&UserCountBySegmentResponse{
		Segment: "sport",
		Count:   10,
	}, nil)

	res, err := sut.GetUserCountBySegmentation(context.Background(), UserCountBySegmentRequest{
		Name: "sport",
	})

	asserts.Nil(err)
	asserts.NotNil(res)
	redisMock.AssertCalled(t, "CountUsersInSegment", mock.Anything, mock.Anything)
	repositoryMock.AssertNotCalled(t, "CountUsersInSegment", mock.Anything, mock.Anything)
}

func Test_GetUserCountBySegmentation_WhenCacheIsNull(t *testing.T) {
	asserts := assert.New(t)

	redisMock := new(RedisMock)
	repositoryMock := new(RepositoryMock)

	sut := ServiceImpl{
		Repository: repositoryMock,
		Redis:      redisMock,
	}

	redisResponse := &UserCountBySegmentResponse{}
	redisMock.On("CountUsersInSegment", mock.Anything, mock.Anything).Return(redisResponse, errors.New("redisErr"))
	repositoryMock.On("CountUsersInSegment", mock.Anything, mock.Anything).Return(&UserCountBySegmentResponse{
		Segment: "sport",
		Count:   10,
	}, nil)
	res, err := sut.GetUserCountBySegmentation(context.Background(), UserCountBySegmentRequest{
		Name: "sport",
	})
	asserts.NotNil(res)
	asserts.Nil(err)
	redisMock.AssertCalled(t, "CountUsersInSegment", mock.Anything, mock.Anything)
	repositoryMock.AssertCalled(t, "CountUsersInSegment", mock.Anything, mock.Anything)
}

func Test_StoreUserSegment_CheckOldCacheDeleted_ok(t *testing.T) {
	asserts := assert.New(t)

	redisMock := new(RedisMock)
	repositoryMock := new(RepositoryMock)

	sut := ServiceImpl{
		Repository: repositoryMock,
		Redis:      redisMock,
	}

	repositoryMock.On("StoreUserInSegment", mock.Anything, mock.Anything).Return(nil)
	redisMock.On("ClearUserCacheInSegment", mock.Anything, mock.Anything)
	redisMock.wg.Add(1)
	err := sut.StoreUserSegment(context.Background(), StoreUserSegmentRequest{"1", "sport"})
	asserts.Nil(err)

	if redisMock.currentCount() != 1 {
		t.Error("redis not called")
	}
}
