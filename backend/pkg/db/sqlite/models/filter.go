package models

type Filter struct {
	GroupID        string `json:"groupID"`
	AuthorID         string `json:"authorID"`
	LikedByUserID    string `json:"likedByUserID"`
	DisLikedByUserID string `json:"disLikedByUserID"`
}

// func (f *Filter) IsCheckedCategory(id string) bool {
// 	for _, c := range f.CategoryID {
// 		if id == c {
// 			return true
// 		}
// 	}
// 	return false
// }
