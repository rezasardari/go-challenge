package pkg

import "context"

type Redis interface {
	CountUsersInSegment(context context.Context, title string) (*UserCountBySegmentResponse, error)
	StoreUserCountInSegment(context context.Context, data UserCountBySegmentResponse)
	ClearUserCacheInSegment(ctx context.Context, title string)
}
