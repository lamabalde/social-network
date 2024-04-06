package webmodel

type FollowingReply struct {
	Id           string `json:"id"`
	FollowStatus string `json:"followStatus"`
}

type FollowResponseWithNotificationID struct {
	Id     int
	UserID string
}
