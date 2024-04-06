import { computed, ref, watch } from "vue";
import ws from "../ws";
import { followersComponent } from "./followerComponent";
import { useRoute } from "vue-router";
import { convertDateToDateString } from "../helpers/helpers";


export const currentUser = ref({
  userInfo: {
    dateCreate: "",
    dateOfBirth: "",
    email: "",
    expirySession: "",
    firstName: "",
    id: "",
    lastName: "",
    sessionUuid: "",
    userName: "",
    profileType: "0",
    aboutMe: "",
  },
  loaded: false,
  addCurrentSession(payload) {
    Object.assign(this.userInfo, payload.data.User);
    this.loaded = true;
  },
  changeProfileType(profileType) {
    ws.request("setProfileType", profileType);
  },
  
});

export const userComponent = ref({
  userInfo: {
    dateCreate: "",
    dateOfBirth: "",
    email: "",
    expirySession: "",
    firstName: "",
    id: "",
    lastName: "",
    sessionUuid: "",
    userName: "",
    profileType: 1,
    aboutMe: "",
    profileImg: "",
  },
  profileLoaded : false,
  checkUserExists(userID) {
    ws.request("checkUserExist", userID);
  },
  renderUserExist(payload) {
    if (!payload.data.exists) {
      window.location.href = "/pageNotFound";
    }
  },
  getUserInfo(userID) {
    ws.request("getUserProfile", userID);
  },
  renderUserInfo(payload) {
    Object.assign(this.userInfo, payload.data);
    this.loaded = true;
  },
  renderProfileType(payload) {
    this.userInfo.profileType = payload.data;
  },
  
});

export const onlineUserTracker = ref({
  onlineUsers: {},
  renderOnlineUsers(payload) {
    if (payload.data != null) {
      payload.data.forEach((user) => {
        this.onlineUsers[user.id] = user.userName;
      });
    }
  },
  renderNewOnlineUser(payload) {
    this.onlineUsers[payload.data.id] = payload.data.userName;
  },
  renderOfflineUser(payload) {
    delete this.onlineUsers[payload.data.id]
  }
});

watch(
  () => userComponent.value.userInfo.id,
  (userID) => {
    if (userID != "") {
      userComponent.value.profileLoaded = true;
      userComponent.value.userInfo.dateOfBirth = convertDateToDateString(
        userComponent.value.userInfo.dateOfBirth
      );
    }
  }
);