package pkg

import (
	"context"
)

type Service interface {
	GetUserCountBySegmentation(context context.Context, req UserCountBySegmentRequest) (*UserCountBySegmentResponse, error)
	StoreUserSegment(ctx context.Context, req StoreUserSegmentRequest) error
	ArchiveExpiredSegment(ctx context.Context) error
	UpdateCacheData(ctx context.Context) error
}

type ServiceImpl struct {
	Repository Repository
	Redis      Redis
}

func (s *ServiceImpl) NewEstimationService(repository Repository, redis Redis) *ServiceImpl {
	return &ServiceImpl{
		Repository: repository,
		Redis:      redis,
	}
}

func (s *ServiceImpl) GetUserCountBySegmentation(context context.Context, req UserCountBySegmentRequest) (*UserCountBySegmentResponse, error) {
	if cache, cErr := s.Redis.CountUsersInSegment(context, req.Name); cErr == nil {
		return cache, nil
	}
	res, err := s.Repository.CountUsersInSegment(context, req)
	if err != nil {
		return nil, err
	}
	go s.Redis.StoreUserCountInSegment(context, *res)
	return res, nil
}

func (s *ServiceImpl) StoreUserSegment(ctx context.Context, req StoreUserSegmentRequest) error {
	if err := s.Repository.StoreUserInSegment(ctx, req); err != nil {
		return err
	}
	go s.Redis.ClearUserCacheInSegment(ctx, req.Segment)
	return nil
}

func (s *ServiceImpl) ArchiveExpiredSegment(ctx context.Context) error {
	if err := s.Repository.ArchiveExpiredData(ctx); err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) retrieveFromCache(ctx context.Context, segment string) (*UserCountBySegmentResponse, error) {
	resp, err := s.Redis.CountUsersInSegment(ctx, segment)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ServiceImpl) UpdateCacheData(ctx context.Context) error {
	resp, err := s.Repository.CountUsersInAllSegments(ctx)
	if err != nil {
		return err
	}
	rErr := s.Redis.StoreUserCountsInAllSegment(ctx, resp)
	if rErr != nil {
		return rErr
	}
	return nil
}

func (s *ServiceImpl) updateCacheData(ctx context.Context) error {
	res, err := s.Repository.CountUsersInAllSegments(ctx)
	if err != nil {
		return err
	}
	rErr := s.Redis.StoreUserCountsInAllSegment(ctx, res)
	if rErr != nil {
		return rErr
	}
	return nil
}
