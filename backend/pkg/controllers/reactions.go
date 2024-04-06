package controllers

import (
	"errors"
	"fmt"

	"social-network/backend/application"
	"social-network/backend/pkg/controllers/liker"
	"social-network/backend/pkg/controllers/wsconnection"
	"social-network/backend/pkg/db/sqlite/models"
	"social-network/backend/pkg/webmodel"
	"social-network/backend/pkg/webmodel/parse"
)

// type reactionData struct {
// 	UserID      int64 `json:"userID"`
// 	PostID      int64 `json:"postID"`
// 	Reaction    bool  `json:"reaction"`    // True for like, false for dislike
// 	AddOrRemove bool  `json:"addOrRemove"` // True to add, false to remove
// 	IsPost      bool  `json:"ispost"`
// }

type Response struct {
	AmountOfReactions int `json:"reactionAmount"` // Likes or dislikes, corresponding on request
}

func ReplyReaction(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		reactionData, err := parse.PayloadToReaction(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload in reaction request: %s", message.Payload), err)
		}

		errmessage := reactionData.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		var react liker.Liker
		switch reactionData.MessageType {
		case models.POST:
			react = liker.NewLikePost(currConnection.Client.UserID, currConnection.Client.UserName, reactionData)
		case models.COMMENT:
			react = liker.NewLikePost(currConnection.Client.UserID, currConnection.Client.UserName, reactionData)
		default:
			return nil, currConnection.WSError(fmt.Sprintf("unexpected message type in reaction request %s", reactionData.MessageType), errors.New("unexpected message type"))
		}

		if err = liker.SetLike(app.DBModel, react, reactionData.Reaction); err != nil {
			return nil, currConnection.WSError("DB error during reaction handling", err)
		}

		// get the new number of likes/dislikes
		newReactions, err := react.GetLikesNumbers(app.DBModel)
		if err != nil {
			return nil, currConnection.WSError("get new reactions failed", err)
		}

		return newReactions, nil
	}
}
