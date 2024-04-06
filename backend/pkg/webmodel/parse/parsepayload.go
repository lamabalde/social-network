package parse

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

func PayloadToInt(payload json.RawMessage) (int, error) {
	var number int
	err := json.Unmarshal(payload, &number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func PayloadToString(payload json.RawMessage) (string, error) {
	var str string
	err := json.Unmarshal(payload, &str)
	if err != nil {
		return "", err
	}

	return str, nil
}

func PayloadToUserCredential(payload json.RawMessage) (webmodel.UserCredentials, error) {
	var uC webmodel.UserCredentials
	decPayload := make(json.RawMessage, len(payload)-2)
	if len(payload) < 2 {
		return uC, errors.New("not enough data")
	}
	n, err := base64.StdEncoding.Decode(decPayload, payload[1:len(payload)-1]) // not to take into account "" in the payload
	if err != nil {
		return uC, err
	}
	err = json.Unmarshal(decPayload[:n], &uC) //:n - to delete empty (0x00) bytes
	return uC, err
}

func PayloadToPost(payload json.RawMessage) (webmodel.Post, error) {
	var post webmodel.Post
	err := json.Unmarshal(payload, &post)
	return post, err
}

func PayloadToComment(payload json.RawMessage) (webmodel.Comment, error) {
	var comment webmodel.Comment
	err := json.Unmarshal(payload, &comment)
	return comment, err
}

func PayloadToChatMessage(payload json.RawMessage) (webmodel.ChatMessage, error) {
	var message webmodel.ChatMessage
	err := json.Unmarshal(payload, &message)
	return message, err
}

func PayloadToReaction(payload json.RawMessage) (webmodel.Reaction, error) {
	var react webmodel.Reaction
	err := json.Unmarshal(payload, &react)
	return react, err
}

func PayloadToGroup(payload json.RawMessage) (webmodel.Group, error) {
	var group webmodel.Group
	err := json.Unmarshal(payload, &group)
	return group, err
}

func PayloadToGroupUser(payload json.RawMessage) (webmodel.GroupUser, error) {
	var userInGroup webmodel.GroupUser
	err := json.Unmarshal(payload, &userInGroup)
	return userInGroup, err
}

func PayloadToPostsInGroupPortion(payload json.RawMessage) (webmodel.PostsInGroupPortion, error) {
	var postsInGroupPortion webmodel.PostsInGroupPortion
	err := json.Unmarshal(payload, &postsInGroupPortion)
	return postsInGroupPortion, err
}

func PayloadToUserPostsPortion(payload json.RawMessage) (webmodel.UserPosts, error) {
	var postsInGroupPortion webmodel.UserPosts
	err := json.Unmarshal(payload, &postsInGroupPortion)
	return postsInGroupPortion, err
}

func PayloadToGroupEvent(payload json.RawMessage) (webmodel.Event, error) {
	var groupEvent webmodel.Event
	err := json.Unmarshal(payload, &groupEvent)
	return groupEvent, err
}

func PayloadToUserOptionForEvent(payload json.RawMessage) (webmodel.UserOptionForEvent, error) {
	var groupEvent webmodel.UserOptionForEvent
	err := json.Unmarshal(payload, &groupEvent)
	return groupEvent, err
}

func PayloadToFollowResponseWithNotificationID(payload json.RawMessage) (webmodel.FollowResponseWithNotificationID, error) {
	var groupEvent webmodel.FollowResponseWithNotificationID
	err := json.Unmarshal(payload, &groupEvent)
	return groupEvent, err
}
