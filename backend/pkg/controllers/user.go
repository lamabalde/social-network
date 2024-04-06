package controllers

import (
	"fmt"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyGetUserProfile(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid userID '%s'", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		user, err := app.DBModel.GetUserByID(userID)
		if err != nil {
			return nil, handleErrGetFromDBByID(app, currConnection, message, "user", userID, err)
		}
		if user == nil {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("cannot find a user with id '%s'", userID))
		}
		user.ProfileImg = getUserImgUrl(userID)
		fmt.Println("profile img: ", user.ProfileImg)
		return user, nil
	}
}

func ReplySetProfileType(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		profileType, err := parse.PayloadToInt(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid userID '%s'", message.Payload), err)
		}

		if profileType < 0 || profileType >= webmodel.NUM_OF_PROFILE_TYPES {
			return nil, currConnection.WSBadRequest(message, "wrong profile type")
		}

		err = app.DBModel.SetUsersProfileType(currConnection.Client.UserID, profileType)
		if err != nil {
			return nil, currConnection.WSError("get user from DB failed", err)
		}

		return profileType, nil
	}
}

func getUserImgUrl(userID string) string { // need to get the file with extension based on the filename
	imgUrl := strings.Join(strings.Split(helpers.FindFile("data/img/profile/", userID+".*"), "/")[1:], "/")
	return imgUrl
}
