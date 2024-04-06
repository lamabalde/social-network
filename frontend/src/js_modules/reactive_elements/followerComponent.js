import ws from "../ws";
import { ref } from "vue";
import { currentUser } from "./userComponent";
import { notificationsComponent } from "./notificationComponent";


export const followersComponent = ref({
  followers: [],
  following: [],
  followStatus: "not following",
  displayFollowBar: "followers",
  getFollowers(userID) {
    ws.request("userFollowers", userID);
  },
  getFollowing(userID) {
    ws.request("userFollowing", userID);
  },
  renderFollowers(payload) {
    this.followers = payload.data; // displays the followers
  },
  renderFollowing(payload) {
    this.following = payload.data; // displays the followers
  },
  followUser(userID) {
    ws.request("followUser", userID);
  },
  unFollowUser(userID) {
    ws.request("unFollowUser", userID);
  },
  getFollowStatus(userID) {
    ws.request("getFollowStatus", userID);
  },
  renderFollowStatus(payload) {
    this.followStatus = payload.data;
  },
  displayFollowingBar() {
    this.displayFollowBar = "following";
  },
  displayFollowerBar() {
    this.displayFollowBar = "followers";
  },
  acceptFollowRequest(payload) {
    ws.request("acceptFollowRequest", payload);
  },
  declineFollowRequest(payload) {
    ws.request("declineFollowRequest", payload);
  },
});
