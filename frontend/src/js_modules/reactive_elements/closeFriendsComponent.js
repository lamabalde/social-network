import { ref } from "vue";
import ws from "../ws";

export const closeFriendsComponent = ref({
    closeFriends: [],
    addCloseFriend(userID) {
        ws.request("addCloseFriend", userID)
    },
    renderCloseFriends(payload) {
        this.closeFriends = payload.data
    }
})