package helpers

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const SEPARATOR = ";"

func GenerateNewUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil // convert 16 byte uuid into string
}

/*
joins slice of strings to sql.NullString, skips empty strings or strings that comprise only whitespaces
*/
func JoinToNullString(slice []string) (nullStr sql.NullString) {
	if len(slice) == 0 {
		return
	}

	var str string
	var i int
	// find 1st not empty string
	for i = 0; i < len(slice); i++ {
		if s := strings.TrimSpace(slice[i]); s != "" {
			str = s
			break
		}
	}
	// add only not empty strings
	for i++; i < len(slice); i++ {
		if s := strings.TrimSpace(slice[i]); s != "" {
			str += SEPARATOR + s
		}
	}

	if len(str) != 0 {
		nullStr.String = str
		nullStr.Valid = true
	}

	return
}

/*
converts string containing a list of images file names to an array
*/
func SplitNullString(imagesStr sql.NullString) []string {
	if imagesStr.Valid {
		imagesNames := strings.Split(imagesStr.String, SEPARATOR)
		for i := 0; i < len(imagesNames)-1; i++ {
			if imagesNames[i] == "" {
				imagesNames = append(imagesNames[:i], imagesNames[i+1])
			}
		}
		if imagesNames[len(imagesNames)-1] == "" {
			imagesNames = imagesNames[:len(imagesNames)-1]
		}
		return imagesNames
	}
	return nil
}

func StrToNullString(str string) sql.NullString {
	var nullstr sql.NullString
	if str != "" {
		nullstr.String = str
		nullstr.Valid = true
	}
	return nullstr
}

func GetPostImgUrl(postID string) string { // need to get the file with extension based on the filename
	imgUrl := strings.Join(strings.Split(findFile("data/img/post/", postID+".*"), "/")[1:], "/")
	return imgUrl
}

func findFile(targetDir string, pattern string) string {

	matches, err := filepath.Glob(targetDir + pattern)
	if err != nil {
		fmt.Println(err)
	}
	if len(matches) != 0 {
		return matches[0]
	}
	return ""
}
