package webmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	CurrentSession = "currentSession"
	ERROR          = "ERROR"
	Register       = "register"

	Login  = "login"
	Logout = "logout"

	PostsPortion        = "postsPortion"
	PostsPortionInGroup = "postsPortionInGroup"
	UserPostsPortion    = "userPostsPortion"
	FullPostAndComments = "fullPostAndComments"
	NewPost             = "newPost"
	DeletePost          = "deletePost"
	NewComment          = "newComment"

	UserJoinedChat         = "userJoinedChat"
	UserQuitChat           = "userQuitChat"
	ChattingUsers          = "chattingUsers"
	SendMessageToOpendChat = "sendMessageToOpendChat"
	InputChatMessage       = "inputChatMessage"
	CloseChat              = "closeChat"
	ChatPortion            = "chatPortion"
	GetChatList            = "getChatList"
	SendMessageToChat      = "sendMessageToChat"

	OnlineUsers = "onlineUsers"
	UserOnline  = "newOnlineUser"
	UserOffline = "offlineUser"

	CreateGroup               = "createGroup"
	GetGroupProfile           = "getGroupProfile"
	LeaveGroup                = "leaveGroup"
	RequestToJoinToGroup      = "requestToJoinToGroup"
	InviteToGroup             = "inviteToGroup"
	AcceptRequestToJoinGroup  = "acceptRequestToJoinGroup"
	DeclineRequestToJoinGroup = "declineRequestToJoinGroup"
	AcceptInvitationToGroup   = "acceptInvitationToGroup"
	DeclineInvitationToGroup  = "declineInvitationToGroup"
	GroupMembers              = "groupMembers"
	UserGroups                = "userGroups"

	CreateGroupEvent         = "createGroupEvent"
	GetGroupEvent            = "getGroupEvent"
	SetUserOptionForEvent    = "setUserOptionForEvent"
	ChangeUserOptionForEvent = "changeUserOptionForEvent"
	GetGroupEventsList       = "getGroupEventsList"

	GetUserProfile = "getUserProfile"

	UserFollowers        = "userFollowers"
	UserFollowing        = "userFollowing"
	FollowUser           = "followUser"
	GetFollowStatus      = "getFollowStatus"
	DeclineFollowRequest = "declineFollowRequest"
	AcceptFollowRequest  = "acceptFollowRequest"
	UnFollowUser         = "unFollowUser"

	NewNotification      = "newNotification"
	GetUserNotifications = "getUserNotifications"

	SearchGroupsUsers     = "searchGroupsUsers"
	SearchUsersNotInGroup = "searchUsersNotInGroup"
	SearchNotCloseFriends = "searchNotCloseFriends"

	SetProfileType = "setProfileType" // PUBLIC = 0,	PRIVATE = 1

	LikePost = "likePost"

	AddCloseFriend = "addCloseFriend"

	// verification
	VerifyGroupView = "verifyGroupView"
)

// --------------------- privacy status --------------------------------
const NUM_OF_PROFILE_TYPES = 2

const (
	PUBLIC = iota
	PRIVATE
	CHOSEN // access for chosen users
)

// --------------------- content's portions --------------------------------
const (
	POSTS_ON_POSTSVIEW    = 10
	CHAT_MESSAGES_PORTION = 10
)

// --------------------- notification's types -------------------------------
const (
	NOTE_FOLLOW_REQUEST              = "follow request"
	NOTE_JOIN_GROUP_REQUEST          = "join group request"
	NOTE_INVITE_TO_GROUP             = "invite to group"
	NOTE_JOIN_GROUP_REQUEST_ACCEPTED = "join group request accepted"
	NOTE_JOIN_GROUP_REQUEST_DECLINED = "join group request declined"
	NOTE_INVITE_TO_GROUP_ACCEPTED    = "invite to group accepted"
	NOTE_INVITE_TO_GROUP_REJECTED    = "invite to group rejected"
	NOTE_NEW_PRIVATE_MESSAGE         = "new chat message"
)

var ErrWarning = errors.New("Warning")

/*
	 message payload's data are
		- string for IDs
		- int for offsets
		- struct from this package if the pauyload is json object
*/
type Payload struct {
	Result string `json:"result"`
	Data   any    `json:"data,omitempty"`
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (m *WSMessage) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("Type: %s | Payload: %s\n", m.Type, m.Payload)
}

func (m *WSMessage) IsAuthentification() bool {
	return strings.HasPrefix(m.Type, "login") || strings.HasPrefix(m.Type, "logout") || strings.HasPrefix(m.Type, "register")
}

func (m *WSMessage) CreateReplyToRequestMessage(result string, data any) (json.RawMessage, error) {
	return CreateJSONMessage(m.Type, result, data)
}

func CreateJSONMessage(messageType string, result string, data any) (json.RawMessage, error) {
	wsMessage, err := createWSMessage(messageType, result, data)
	if err != nil {
		return nil, fmt.Errorf("CreateJSONMessage failed: %v", err)
	}

	jsonMessage, err := json.Marshal(wsMessage)
	if err != nil {
		return nil, fmt.Errorf("CreateJSONMessage failed: %v", err)
	}
	return jsonMessage, nil
}

func createWSMessage(messageType string, result string, data any) (WSMessage, error) {
	payload, err := json.Marshal(Payload{Result: result, Data: data})
	if err != nil {
		return WSMessage{}, fmt.Errorf("createWSMessage failed: %v", err)
	}
	message := WSMessage{
		Type:    messageType,
		Payload: payload,
	}
	return message, nil
}

func IsEmpty(field string) bool {
	return strings.TrimSpace(field) == "" || field == "undefined"
}
