import { ref } from "vue";
import ws from "../ws";
import Socket from "../models/webSocketModel";
import { chatMessage } from "../models/chatMessageModel";
import { sortObjArrayByParameter } from "../helpers/helpers";
import { currentUser } from "./userComponent";

export const chatSockets = {} // chatRoomID is the key and the ws connection is the value for every property
export const chatMessageInput = ref("")

export const chatComponent = ref({
  chatList: [], //list of users in the friends list
  messages: [],
  currentOpenChat: {},
  isChatOpen: false,
  messagesLoaded: false,
  offset: 0,
  openCloseChatWindow(recipient) {
    this.isChatOpen = !this.isChatOpen;
    this.offset = 0;
    this.messages = [];
    if (this.isChatOpen) {
      this.chatList.forEach((chatElem, index) => { // set the amount of new messages to 0
        if (chatElem.genericID == recipient.genericID) {
          this.chatList[index].unreadMessages = 0;
        }
      });
      if (recipient.type == 0) {
        this.openPrivateChat(recipient.id);
      } else if (recipient.type == 1) {
        this.openGroupChat(recipient.id);
      }
      this.currentOpenChat = recipient;
    } else {
      const currentSocket = chatSockets[this.currentOpenChat.id];
      currentSocket.closeWebsocket();
      delete chatSockets[this.currentOpenChat.id];

      this.currentOpenChat = {};
      this.messagesLoaded = false;
    }
  },

  openPrivateChat(recipientID) {
    const socket = new Socket();
    chatSockets[recipientID] = socket;
    socket.launchWebsocket(`openChat/private?id=${recipientID}`);

    socket.loadPageAfterWsOpen(function () {
      chatComponent.value.getChatMessages(recipientID);
    });
  },
  openGroupChat(groupID) {
    const socket = new Socket();
    chatSockets[groupID] = socket;
    socket.launchWebsocket(`openChat/group?id=${groupID}`);

    socket.loadPageAfterWsOpen(function () {
      chatComponent.value.getChatMessages(groupID);
    });
  },
  sendMessage(recipientID, message) {
    const socket = chatSockets[recipientID];
    const payload = new chatMessage(message);
    socket.request("sendMessageToChat", payload);
  },
  getChatList(userID) {
    //directed to ws_response_router.js
    ws.request("getChatList", userID);
  },
  renderChatList(payload) {
    //receives data from ReplyGetChatList

    //
    payload.data.forEach((chatElem, index) => {
      if (this.chatList[index]) { 
        // can't do the chatlist sorting within the browser because we don't have last message data, also we'd need to update it with every new message and then sort
        chatElem["unreadMessages"] = this.chatList[index].unreadMessages; // when reordering the chatlist, it keeps the old unreadMessage counter
      } else {
        chatElem["unreadMessages"] = 0;
      }
      
    })
    this.chatList = payload.data;
  },
  getChatMessages(recipientid) {
    const socket = chatSockets[recipientid];
    socket.request("chatPortion", chatComponent.value.offset);
  },
  renderChatMessages(payload) {
    if (payload.data.messages) {
      this.messages = payload.data.messages.concat(this.messages);
      //this.messages = sortObjArrayByParameter(this.messages, "dateCreate");
      setTimeout(function() {
        chatComponent.value.messagesLoaded = true; // for scrolling to the bottom of new messages
      }, 10)
      
    }
  },
  handleNewMessage(payload) {
    
    this.messages.push(payload.data);
      
    this.chatList.forEach((chatElem, index) => {
      if (
        payload.data.genericID == chatElem.id &&
        payload.data.genericID !== this.currentOpenChat.id
      ) {
        // if the incoming message chat is already open, don't show new message counter
        this.chatList[index].unreadMessages += 1;
      }
    })
    setTimeout(function() {
      chatComponent.value.messagesLoaded = true; // for scrolling to the bottom of new messages
    }, 10)
    //this.getChatList(currentUser.value.userInfo.id); // TODO: get it working with sorting user list by last send messages
  }
});
