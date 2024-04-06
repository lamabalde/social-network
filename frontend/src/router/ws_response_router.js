import { useRoute } from "vue-router";
import { commentsComponent } from "../js_modules/reactive_elements/commentsComponent"
import { followersComponent } from "../js_modules/reactive_elements/followerComponent";
import { groupInfoComponent, groupListComponent, groupMembersComponent } from "../js_modules/reactive_elements/groupsComponent"
import { notificationsComponent } from "../js_modules/reactive_elements/notificationComponent";
import { postsComponent } from "../js_modules/reactive_elements/postComponent"
import { searchComponent } from "../js_modules/reactive_elements/searchComponent";
import {
  currentUser,
  onlineUserTracker,
  userComponent,
} from "../js_modules/reactive_elements/userComponent";
import { verifyComponent } from "../js_modules/reactive_elements/verifyPageView";
import { eventInfo, eventsComponent } from "../js_modules/reactive_elements/eventsComponent";
import ws from "../js_modules/ws";
import { findIndexInArray } from "../js_modules/helpers/helpers";
import { chatComponent } from "../js_modules/reactive_elements/chatComponent";

// routes the response from websocket to various functions for components

export const wsResponseRouter = {
  // CURRENTLY WORKING
  // online/offline tracker
  newOnlineUser(payload) {
    onlineUserTracker.value.renderNewOnlineUser(payload);
  },
  onlineUsers(payload) {
    onlineUserTracker.value.renderOnlineUsers(payload);
  },
  offlineUser(payload) {
    onlineUserTracker.value.renderOfflineUser(payload);
  },
  // chat
  inputChatMessage(payload) {
    chatComponent.value.handleNewMessage(payload);
  },
  sendMessageToChat(payload) {
    // irrelevant
  },
  chatPortion(payload) {
    chatComponent.value.renderChatMessages(payload);
  },
  chattingUsers(payload) {},
  userJoinedChat(payload) {},
  userQuitChat(payload) {},
  getChatList(payload) {
    chatComponent.value.renderChatList(payload);
  },
  // search
  searchGroupsUsers(payload) {
    searchComponent.value.renderSearchResults(payload);
  },
  searchUsersNotInGroup(payload) {
    searchComponent.value.renderSearchResults(payload);
  },
  searchNotCloseFriends(payload) {
    searchComponent.value.renderCloseFriendSearchResults(payload);
  },
  // profile page
  getUserProfile(payload) {
    userComponent.value.renderUserInfo(payload);
  },
  //posts, comments
  newPost(payload) {
    postsComponent.value.renderNewPost(payload);
  },
  postsPortion(payload) {
    postsComponent.value.renderPostsPortion(payload);
  },
  postsPortionInGroup(payload) {
    postsComponent.value.renderPostsPortion(payload);
  },
  userPostsPortion(payload) {
    postsComponent.value.renderPostsPortion(payload);
  },
  deletePost(payload) {
    postsComponent.value.deletePost(payload);
  },
  fullPostAndComments(payload) {
    commentsComponent.value.renderAllComments(payload);
  },
  newComment(payload) {
    commentsComponent.value.renderNewComment(payload);
  },
  //likes
  likePost(payload) {
    postsComponent.value.renderPostLike(payload);
  },
  //session
  currentSession(payload) {
    currentUser.value.addCurrentSession(payload);
  },
  //groups
  groupMembers(payload) {
    groupMembersComponent.value.renderMembers(payload);
  },
  userGroups(payload) {
    groupListComponent.value.renderGroups(payload);
  },
  leaveGroup(payload) {
    groupMembersComponent.value.renderLeaveGroup();
  },
  getGroupProfile(payload) {
    groupInfoComponent.value.renderGroupInfo(payload);
  },
  inviteToGroup(payload) {
    console.log(payload);
    const index = searchComponent.value.searchResults.findIndex((i) => {
      return i.id === payload.data;
    });
    searchComponent.value.searchResults.splice(index, 1);
  },
  requestToJoinToGroup(payload) {
    verifyComponent.value.verifyGroupView(payload.data);
  },
  acceptInvitationToGroup(payload) {
    notificationsComponent.value.getNotifications(
      currentUser.value.userInfo.id
    );
    verifyComponent.value.verifyGroupView(payload.data);
  },
  //followers
  userFollowers(payload) {
    followersComponent.value.renderFollowers(payload);
  },
  userFollowing(payload) {
    followersComponent.value.renderFollowing(payload);
  },
  getFollowStatus(payload) {
    followersComponent.value.renderFollowStatus(payload);
  },
  followUser(payload) {
    if (payload.data.followStatus == "following") {
      followersComponent.value.followStatus = payload.data.followStatus;
      followersComponent.value.getFollowers(payload.data.id); // to re-render follower/following list
      followersComponent.value.getFollowing(payload.data.id);
    } else {
      followersComponent.value.followStatus = payload.data;
    }
  },
  unFollowUser(payload) {
    followersComponent.value.followStatus = payload.data;
    followersComponent.value.getFollowers(payload.data);
    followersComponent.value.getFollowing(payload.data);
  },
  declineFollowRequest(payload) {
    notificationsComponent.value.removeNotification(payload.data.Id);
  },
  acceptFollowRequest(payload) {
    notificationsComponent.value.removeNotification(payload.data.Id);
  },
  addCloseFriend(payload) {
    const index = findIndexInArray(
      searchComponent.value.closeFriendSearchResults,
      payload.data,
      "id"
    );
    searchComponent.value.closeFriendSearchResults.splice(index, 1);
  },
  //misc
  setProfileType(payload) {
    userComponent.value.renderProfileType(payload);
  },
  verifyGroupView(payload) {
    verifyComponent.value.handleGroupVerification(payload);
  },
  // events
  getGroupEventsList(payload) {
    eventsComponent.value.eventsList = payload.data;
  },
  getGroupEvent(payload) {
    eventInfo.value.renderEventInfo(payload);
  },
  createGroupEvent(payload) {
    eventsComponent.value.renderNewGroupEvent(payload);
  },
  setUserOptionForEvent(payload) {
    eventsComponent.value.getEventInfo(payload.data.eventID);
  },
  // notifications
  getUserNotifications(payload) {
    notificationsComponent.value.renderNotifications(payload);
  },
  newNotification(payload) {
    notificationsComponent.value.newNotification(payload);
  },
  // WIP

  declineInvitationToGroup(payload) {
    notificationsComponent.value.getNotifications(
      currentUser.value.userInfo.id
    );
  },
  acceptRequestToJoinGroup(payload) {
    console.log("p");
    notificationsComponent.value.getNotifications(
      currentUser.value.userInfo.id
    );
  },
  declineRequestToJoinGroup(payload) {
    console.log("p");
    notificationsComponent.value.getNotifications(
      currentUser.value.userInfo.id
    );
  },
  createGroup(payload) {
    groupListComponent.value.renderNewGroup(payload);
  },
  addUserToGroup(payload) {
    // TODO: make it work with requesting/accepting
    window.location.reload();
  },

  // TODO

  sendMessageToOpendChat(payload) {
    "sends a message from another user to current client";
  },

  goingToEvent(payload) {},
  notGoingToEvent(payload) {},
  getEventInfo(payload) {},
};