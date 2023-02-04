package pg

import (
	"context"
	"estimation-service/pkg"
	"gorm.io/gorm"
	"time"
)

type RepositoryImpl struct {
	DB *gorm.DB
}

func (r *RepositoryImpl) CountUsersInSegment(ctx context.Context, req pkg.UserCountBySegmentRequest) (*pkg.UserCountBySegmentResponse, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(UserSegmentModel{}).Count(&count).Where("segment = ? AND created_at >= NOW() - INTERVAL '2 weeks'", req.Name).Error
	if err != nil {
		return nil, err
	}
	return &pkg.UserCountBySegmentResponse{
		Segmentation: req.Name,
		Count:        int(count),
	}, nil
}

func (r *RepositoryImpl) StoreUserInSegment(ctx context.Context, req pkg.StoreUserSegmentRequest) error {
	userSegment := &UserSegmentModel{
		UserID:  req.UserID,
		Segment: req.Segment,
	}
	if err := r.DB.WithContext(ctx).Create(userSegment).Error; err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) ArchiveExpiredData(ctx context.Context) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var records []UserSegmentModel
		twoWeeksAgo := time.Now().AddDate(0, 0, -14)
		if err := tx.Where("created_at < ?", twoWeeksAgo).Find(&records).Error; err != nil {
			return err
		}
		archivedUserSegment := toUserSegmentArchivedModel(records)
		if err := tx.Create(archivedUserSegment).Error; err != nil {
			return err
		}

		if err := tx.Delete(records).Error; err != nil {
			return err
		}
		return nil
	})
}
