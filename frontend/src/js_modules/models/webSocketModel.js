import { ref } from "vue";
import {payloadModel} from "../models/payloadModel"
import { wsResponseRouter } from "../../router/ws_response_router";

export default class Socket {
  constructor() {
    this.connection = new WebSocket("ws://temp"); // to prevent some errors in components
    this.loaded = ref(false);
  }
  launchWebsocket(url) {
    this.connection = new WebSocket(`ws://localhost:8000/${url}`);

    this.connection.onopen = function (event) {
      console.log("WebSocket connection opened:", event);
      this.loaded = true;
    };

    this.connection.onmessage = this.handleMessages;

    this.connection.onclose = function (event) {
      console.log("WebSocket connection closed:", event);
    };
    this.connection.onerror = function (event) {
      console.log("WebSocket error:", event);
    };
  }
  handleMessages = (event) => {
    //Split by delimiter and remove any empty string in array with filter
    const rawMessages = event.data.split("\n").filter(Boolean);
    // Parse each message
    rawMessages.forEach((rawMessage) => {
      try {
        const message = JSON.parse(rawMessage);
        this.emit(message.type, message.payload);
      } catch (error) {
        console.error("Unable to parse JSON string:", error);
      }
    });
  };
  emit(event, data) {
    //console.log("Type: ", event, " Payload:", data); // for troubleshooting
    if (data.result == "success") {
      wsResponseRouter[event](data); // routes the data to a handler based on the event
    } else {
      if (event == "getUserProfile") {
        window.location.href = "/userNotFound";
      }
      throw new Error("Failed to complete operation: ", data);
    }
  }
  closeWebsocket() {
    this.connection.close();
  }
  request(type, payload) {
    // handles the websocket request message and converts to json as well
    const init = JSON.stringify(new payloadModel(type, payload));
    this.connection.send(init);
  }
  loadPageAfterWsOpen(loadPage) {
    if (this.connection.readyState == 1) {
      loadPage();
    } else {
      this.connection.addEventListener("open", (e) => {
        loadPage();
        e.target.removeEventListener("open", e);
      });
    }
  }
}
