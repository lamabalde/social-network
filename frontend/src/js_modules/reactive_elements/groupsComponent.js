
import { ref } from "vue";
import ws from "../ws";
import { formToObj } from "../helpers/helpers";

export const groupListComponent = ref({
  groups: [],
  getGroups(userID) {
    ws.request("userGroups", userID);
  },
  renderGroups(payload) {
    this.groups = payload.data; // displays the groups
  },
  async createGroup(formEvent) {
    const formObj = await formToObj(formEvent);
    ws.request("createGroup", formObj);
  },
  renderNewGroup(payload) {
    this.groups.push(payload.data)
  }
});

export const groupMembersComponent = ref({
  members: [],
  currentMemberStatus: -1,
  getMembers(groupID) {
    ws.request("groupMembers", groupID);
  },
  renderMembers(payload) {
    this.members = payload.data; // displays the followers
  },

  renderJoinGroup(payload) {
    window.location.reload();
  },
  leaveGroup(groupID) {
    ws.request("leaveGroup", groupID);
  },
  renderLeaveGroup(groupID) {
    window.location.reload();
  },
  inviteUserToGroup(groupID, userID) {
    const payload = {
      groupID: groupID,
      user: userID,
    };
    ws.request("inviteToGroup", payload);
  },
  joinGroup(groupID) {
    ws.request("requestToJoinToGroup", groupID);
  },
  acceptGroupInvite(notificationID) {
    ws.request("acceptInvitationToGroup", notificationID);
  },
  declineGroupInvite(notificationID) {
    ws.request("declineInvitationToGroup", notificationID);
  },
  acceptGroupJoinRequest(notificationID) {
    ws.request("acceptRequestToJoinGroup", notificationID);
  },
  declineGroupJoinRequest(notificationID) {
    ws.request("declineRequestToJoinGroup", notificationID);
  },
});

export const groupInfoComponent = ref({
  // for group info and page verification
  groupInfo: {
    name: "",
    dateCreate: "",
    id: "",
    description: "",
    creator_id: "451ed6cf-1275-4f1b-b15c-897129eb61ee",
  },
  getGroupInfo(groupID) {
    ws.request("getGroupProfile", groupID);
  },
  renderGroupInfo(payload) {
    Object.assign(this.groupInfo, payload.data);
  },
});