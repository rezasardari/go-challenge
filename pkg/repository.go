package pkg

import "context"

type Repository interface {
	CountUsersInSegment(ctx context.Context, req UserCountBySegmentRequest) (*UserCountBySegmentResponse, error)
	StoreUserInSegment(ctx context.Context, req StoreUserSegmentRequest) error
	ArchiveExpiredData(ctx context.Context) error
	CountUsersInAllSegments(ctx context.Context) ([]UserCountBySegmentResponse, error)
}
