package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ConvertBase64ToImg(encodedImg string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(encodedImg)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func CreateUserImgFromBytes(bytes []byte, userID, imgType string) (string, error) {

	mimeType := http.DetectContentType(bytes)
	fileExtension := strings.Split(mimeType, "/")
	if fileExtension[0] != "image" {
		return "", errors.New("invalid file type as profile img")
	}
	fileName := userID + "." + fileExtension[1]
	f, err := os.Create("data/img/profile/" + fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}
	return fileName, nil
}

func CreateImgFromBytes(bytes []byte, ID, imgType, contentType string) (string, error) {

	mimeType := http.DetectContentType(bytes)
	fileExtension := strings.Split(mimeType, "/")
	if fileExtension[0] != "image" {
		return "", errors.New("invalid file type as img")
	}
	fileName := ID + "." + fileExtension[1]
	f, err := os.Create(fmt.Sprintf("data/img/%s/%s", contentType, fileName))
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}
	return fileName, nil
}
