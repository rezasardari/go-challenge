package pkg

type UserCountBySegmentRequest struct {
	Name string
}

type UserCountBySegmentResponse struct {
	Segment string
	Count   int
}

type StoreUserSegmentRequest struct {
	UserID  string
	Segment string
}
