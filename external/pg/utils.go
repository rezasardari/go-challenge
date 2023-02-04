package pg

func toUserSegmentArchivedModel(req []UserSegmentModel) []UserSegmentArchivedModel {
	var userSegmentsArchived []UserSegmentArchivedModel
	for _, userSegment := range req {
		userSegmentsArchived = append(userSegmentsArchived, UserSegmentArchivedModel{
			UserID:  userSegment.UserID,
			Segment: userSegment.Segment,
		})
	}
	return userSegmentsArchived
}
