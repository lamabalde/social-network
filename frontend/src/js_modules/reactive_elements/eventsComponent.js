import { ref } from "vue";
import { convertDateToLocale, formToObj } from "../helpers/helpers";
import ws from "../ws";

export const eventsComponent = ref({
  eventsList: [],
  getGroupEvents(groupID) {
    ws.request("getGroupEventsList", groupID);
  },
  renderGroupEvents(payload) {
    this.eventsList = payload.data;
  },
  async createGroupEvent(e) {
    const payload = await formToObj(e);
    const formattedDate = new Date(payload.dateEvent);
    payload.dateEvent = formattedDate
    // turn the event date into RFC 3339
    ws.request("createGroupEvent", payload);
  },
  goingToEvent(eventID) {
    const init = {eventID: eventID, option: 1}
    ws.request("setUserOptionForEvent", init);
  },
  noyGoingToEvent(eventID) {
    const init = { eventID: eventID, option: 0 };
    ws.request("setUserOptionForEvent", init);
  },
  getEventInfo(eventID) {
    ws.request("getGroupEvent", eventID);
  },
  renderNewGroupEvent(payload) {
    if (this.eventsList) {
      this.eventsList.push(payload.data);
    } else {
      this.eventsList = [payload.data]
    }
  }
});

export const eventInfo = ref({
  info: {
    id: "string",
    title: "string",
    description: "string",
    dateCreate: "time.Time",
    dateEvent: "time.Time",
    groupID: "string",
    creatorID: "string",
    creatorName: "string",
    userOptions: [],
  },
  going: [],
  notGoing: [],
  renderEventInfo(payload) {
    this.going = []
    this.notGoing = []
    if (payload.data.userOptions != null) {
      payload.data.userOptions.forEach((user) => {
        if (user.option == 1) {
          this.going.push(user);
        } else {
          this.notGoing.push(user);
        }
      });
    }
    payload.data.dateEvent = convertDateToLocale(payload.data.dateEvent);
    Object.assign(this.info, payload.data);
  },
});
/* 
{
      option: 0,
      userID: "string",
      userName: "string",
    } */