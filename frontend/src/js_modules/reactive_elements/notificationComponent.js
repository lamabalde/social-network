
import { ref } from "vue";
import ws from "../ws";
import { verifyComponent } from "./verifyPageView";

export const notificationsComponent = ref({
  notifications: [],
  getNotifications(userID) {
    ws.request("getUserNotifications", userID);
  },
  renderNotifications(payload) {
    if (payload.data!= null) {
      this.notifications = payload.data;
    } else {
      this.notifications = []
    }
    
  },
  newNotification(payload) {
    // add notification visual element

    this.notifications.unshift(payload.data);
    handleNotificationTypes(payload.data)
  },
  removeNotification(notificationID) {
    let i = this.notifications
      .map((notification) => notification.id)
      .indexOf(notificationID);
    this.notifications.splice(i, 1);
  }
});

function handleNotificationTypes(notification) {
  switch (type) {
    case "join group request accepted":
        verifyComponent.value.verifyGroupView(notification.groupId.String)
      break;

    default:
      break;
  }
}