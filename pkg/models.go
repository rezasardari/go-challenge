package pkg

type UserCountBySegmentRequest struct {
	Name string `form:"name" json:"name"`
}

type UserCountBySegmentResponse struct {
	Segment string `json:"segment"`
	Count   int    `json:"count"`
}

type StoreUserSegmentRequest struct {
	UserID  string `json:"user_id"`
	Segment string `json:"segment"`
}
