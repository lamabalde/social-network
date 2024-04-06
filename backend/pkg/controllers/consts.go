package controllers

// form fields
const (
	F_NAME         = "name"
	F_PASSWORD     = "password"
	F_EMAIL        = "email"
	F_CONTENT      = "text"
	F_IMAGES       = "images"
	F_AUTHORID     = "authorID"
	F_THEME        = "title"
	F_CATEGORIESID = "categoriesID"
	F_LIKEBY       = "likedby"
	F_DISLIKEBY    = "dislikedby"
)

const POST_PREVIEW_LENGTH = 450


const USER_IMAGES_DIR = "./images"

const (
	MaxFileUploadSize = 20 << 20               // 20MB
	MaxUploadSize     = 10 * MaxFileUploadSize // 10 files by 20MB
)