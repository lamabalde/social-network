import { ref } from "vue";

export const eventsRender = ref({
    viewEventForm: false,
    viewEventsList: false,
    viewEventGoing: "going",
    displayEventGoing() {
        this.viewEventGoing = "going"
    },
    displayEventNotGoing() {
        this.viewEventGoing = "notGoing"
    }
})

export const groupPageRender = ref({
  viewGroupAboutList: false,
  viewAdminList: false
});