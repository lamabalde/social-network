package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

// for handling comment image with ajax

func SubmitCommentImg(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		imageData := &webmodel.CommentImg{}
		json.NewDecoder(r.Body).Decode(&imageData)

		if imageData.Image != "" {
			img := strings.Split(imageData.Image, ",")
			bytes, _ := helpers.ConvertBase64ToImg(img[1])
			fileName, _ := helpers.CreateImgFromBytes(bytes, imageData.CommentID, img[0], "comment")
			app.DBModel.ModifyComment(imageData.CommentID, "", []string{fileName})
			payload, _ := json.Marshal(fileName)
			w.WriteHeader(http.StatusOK)
			w.Write(payload)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
