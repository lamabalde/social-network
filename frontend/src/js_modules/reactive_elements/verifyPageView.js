import { computed, ref } from "vue";
import ws from "../ws";
import { useRoute } from "vue-router";
import { followersComponent } from "./followerComponent";
import { currentUser, userComponent } from "./userComponent";
import { groupMembersComponent } from "./groupsComponent";

export const verifyComponent = ref({
  verified: false,
  finishedVerification: false,
  rejectUser() {
    this.finishedVerification = true;
    this.verified = false;
  },
  verifyUser() {
    this.finishedVerification = true;
    this.verified = true;
  },
  handleGroupVerification(payload) {
    if (payload.data>=0) {
      this.verifyUser();
      groupMembersComponent.value.currentMemberStatus = payload.data
    } else {
      this.rejectUser();
    }
  },
  verifyGroupView(groupID) {
    this.verified = false;
    ws.request("verifyGroupView", groupID);
  },
  verifyFollowView: () =>
    computed(function () {
      const route = useRoute().params.id;
      // to check if the user is eligible to see the posts
      if (
        followersComponent.value.followStatus == "following" ||
        route == currentUser.value.userInfo.id ||
        userComponent.value.userInfo.profileType == 0
      ) {
        return true;
      }
      return false;
    }),
});