package pg

import "gorm.io/gorm"

type UserSegmentModel struct {
	gorm.Model
	UserID  string
	Segment string `gorm:"index:,option:CONCURRENTLY"`
}

func (UserSegmentModel) TableName() string {
	return "user-segment"
}

type UserSegmentArchivedModel struct {
	gorm.Model
	UserID  string
	Segment string
}

func (UserSegmentArchivedModel) TableName() string {
	return "archived-user-segment"
}
