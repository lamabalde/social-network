package routes

import (
	"errors"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

const (
	WS_REPLYERS_UNI = iota
	WS_REPLYERS_CHAT
)

func CreateUniWsRoutes(app *application.Application) wsconnection.WSmux { // wsconnection.WSmux is map[string]Replier
	wsServer := wsconnection.WSmux{
		InfoLog: app.InfoLog,
		ErrLog:  app.ErrLog,
		Hub:     app.Hub,
	}
	wsServer.WShandlers = map[string]wsconnection.Replier{
		webmodel.UserOffline: sendReplyForLoggedUser(app, wsconnection.FuncReplier(controllers.SendUserStatusToUsers(app, webmodel.UserOffline))),

		webmodel.PostsPortion:        sendReplyForLoggedUser(app, controllers.ReplyPosts(app)),
		webmodel.PostsPortionInGroup: sendReplyForLoggedUser(app, controllers.ReplyPostsInGroup(app)),
		webmodel.UserPostsPortion:    sendReplyForLoggedUser(app, controllers.ReplyUserPostsPortion(app)), // WIP
		webmodel.FullPostAndComments: sendReplyForLoggedUser(app, controllers.ReplyFullPostAndComments(app)),
		webmodel.NewPost:             sendReplyForLoggedUser(app, controllers.ReplyNewPost(app)),
		webmodel.DeletePost:          sendReplyForLoggedUser(app, controllers.ReplyDeletePost(app)),
		webmodel.NewComment:          sendReplyForLoggedUser(app, controllers.ReplyNewComment(app)),

		// webmodel.OpenChat:               sendReplyForLoggedUser(app,controllers.ReplyOpenPrivateChat)),
		webmodel.GetChatList: sendReplyForLoggedUser(app, controllers.ReplyGetChatList(app)),
		// webmodel.CloseChat:              sendReplyForLoggedUser(app,controllers.ReplyCloseChat)),

		webmodel.CreateGroup:               sendReplyForLoggedUser(app, controllers.ReplyCreateGroup(app)),
		webmodel.GetGroupProfile:           sendReplyForLoggedUser(app, controllers.ReplyGetGroupProfile(app)),
		webmodel.LeaveGroup:                sendReplyForLoggedUser(app, controllers.ReplyLeaveGroup(app)),
		webmodel.RequestToJoinToGroup:      sendReplyForLoggedUser(app, controllers.ReplyRequestToJoinToGroup(app)),
		webmodel.InviteToGroup:             sendReplyForLoggedUser(app, controllers.ReplyInviteToGroup(app)),
		webmodel.AcceptRequestToJoinGroup:  sendReplyForLoggedUser(app, controllers.ReplyAcceptRequestToJoinGroup(app)),
		webmodel.DeclineRequestToJoinGroup: sendReplyForLoggedUser(app, controllers.ReplyDeclineRequestToJoinGroup(app)),
		webmodel.AcceptInvitationToGroup:   sendReplyForLoggedUser(app, controllers.ReplyAcceptInvitationToGroup(app)),
		webmodel.DeclineInvitationToGroup:  sendReplyForLoggedUser(app, controllers.ReplyDeclineInvitationToGroup(app)),
		webmodel.GroupMembers:              sendReplyForLoggedUser(app, controllers.ReplyGroupMembers(app)),
		webmodel.UserGroups:                sendReplyForLoggedUser(app, controllers.ReplyUserGroups(app)),

		webmodel.CreateGroupEvent:         sendReplyForLoggedUser(app, controllers.ReplyCreateGroupEvent(app)),
		webmodel.GetGroupEvent:            sendReplyForLoggedUser(app, controllers.ReplyGetGroupEvent(app)),
		webmodel.GetGroupEventsList:       sendReplyForLoggedUser(app, controllers.ReplyGetGroupEventsList(app)),
		webmodel.SetUserOptionForEvent:    sendReplyForLoggedUser(app, controllers.ReplySetUserOptionForEvent(app)),
		webmodel.ChangeUserOptionForEvent: sendReplyForLoggedUser(app, controllers.ReplyChangeUserOptionForEvent(app)),

		webmodel.GetUserProfile: sendReplyForLoggedUser(app, controllers.ReplyGetUserProfile(app)),
		webmodel.SetProfileType: sendReplyForLoggedUser(app, controllers.ReplySetProfileType(app)),

		webmodel.UserFollowers:        sendReplyForLoggedUser(app, controllers.ReplyUserFollowers(app)),
		webmodel.UserFollowing:        sendReplyForLoggedUser(app, controllers.ReplyUserFollowing(app)),
		webmodel.FollowUser:           sendReplyForLoggedUser(app, controllers.ReplyFollowUser(app)),
		webmodel.GetFollowStatus:      sendReplyForLoggedUser(app, controllers.ReplyGetFollowStatus(app)),
		webmodel.AcceptFollowRequest:  sendReplyForLoggedUser(app, controllers.ReplyAcceptFollowRequest(app)),
		webmodel.DeclineFollowRequest: sendReplyForLoggedUser(app, controllers.ReplyDeclineFollowRequest(app)),
		webmodel.UnFollowUser:         sendReplyForLoggedUser(app, controllers.ReplyUnFollowUser(app)),
		webmodel.AddCloseFriend:       sendReplyForLoggedUser(app, controllers.ReplyAddCloseFriend(app)),

		webmodel.SearchGroupsUsers:     sendReplyForLoggedUser(app, controllers.ReplySearchGroupsUsers(app)),
		webmodel.SearchUsersNotInGroup: sendReplyForLoggedUser(app, controllers.ReplySearchUsersNotInGroup(app)),
		webmodel.SearchNotCloseFriends: sendReplyForLoggedUser(app, controllers.ReplySearchUsersNotFriends(app)),

		webmodel.LikePost: sendReplyForLoggedUser(app, controllers.ReplyReaction(app)),

		webmodel.VerifyGroupView: sendReplyForLoggedUser(app, controllers.ReplyVerifyGroupView(app)),

		webmodel.GetUserNotifications: sendReplyForLoggedUser(app, controllers.ReplyGetUserNotifications(app)),

		// needed WIP routes

	}

	return wsServer
}

func CreateChatWsRoutes(app *application.Application) wsconnection.WSmux { // wsconnection.WSmux is map[string]Replier
	wsServer := wsconnection.WSmux{
		InfoLog: app.InfoLog,
		ErrLog:  app.ErrLog,
		Hub:     app.Hub,
	}
	wsServer.WShandlers = map[string]wsconnection.Replier{
		webmodel.SendMessageToChat: sendReplyForLoggedUser(app, controllers.ReplySendMessageToChat(app)),
		webmodel.ChatPortion:       sendReplyForLoggedUser(app, controllers.ReplyChatPortion(app)),
		webmodel.UserOffline:       sendReplyForLoggedUser(app, wsconnection.FuncReplier(controllers.SendUserStatusToChatMembers(app, webmodel.UserQuitChat))),
	}

	return wsServer
}

var ErrNoPost = errors.New("could not find post")

func sendReplyForLoggedUser(app *application.Application, next wsconnection.Replier) wsconnection.Replier {
	return wsconnection.FuncReplier(func(currConnection *wsconnection.UsersConnection, wsMessage webmodel.WSMessage) error {
		err := checkLoggedStatus(app, currConnection, wsMessage)
		if err != nil {
			return err
		}

		return next.SendReply(currConnection, wsMessage)
	})
}

/*
checks the session status, changes it if necessary. Returns nil only if the session has Loggedin status,
otherwise sends an error message to the websocket connection ('conn') and returns an error
*/
func checkLoggedStatus(app *application.Application, currConnection *wsconnection.UsersConnection, requestMessage webmodel.WSMessage) error {
	_, err := currConnection.Session.Tidy(app)
	if err != nil {
		return currConnection.WSError("invalid session status", err)
	}
	if !currConnection.Session.IsLoggedin() {
		return currConnection.WSBadRequest(requestMessage, "not logged in")
	}

	return nil
}
