<script setup>

import Navbar from '../components/navbars/Navbar.vue';
import ws from '../js_modules/ws';
import LeftEventPageInfoBar from '../components/navbars/LeftEventPageInfoBar.vue';
import { eventInfo } from '../js_modules/reactive_elements/eventsComponent';
import { eventsRender } from '../js_modules/reactive_elements/external_view_refs/viewRender';
import { eventsComponent } from '../js_modules/reactive_elements/eventsComponent';
import { useRoute } from 'vue-router';
import ChatComponent from "../components/chat/ChatComponent.vue";
const eventID = useRoute().params.id
function loadPageData() {
    eventsComponent.value.getEventInfo(eventID)
}


/*function initialVerify() { // verify the user is in the group
    verifyComponent.value.verifyGroupView(groupID) // get group id from the event
  }
   watch(()=> verifyComponent.value.finishedVerification, verify => {
    if (verifyComponent.value.verified) {
      ws.loadPageAfterWsOpen(loadPageData)
    }
  }) */

  ws.loadPageAfterWsOpen(loadPageData)
</script>

<template>
    <div id="mainPage" class="content-center flex flex-col bg-[#faf7f2] h-screen">

    <Navbar></Navbar>
    <LeftEventPageInfoBar/>
    <div class="bg-[#faf7f2] mt-16">
    
      
      <div id="allPostData" class="flex flex-col centerContent p-16 border border-gray-200 rounded bg-[#ffffff]">
        <h1 class="flex justify-center items-center font-bold mb-8 text-3xl font-merriweather">Title: {{ eventInfo.info.title }}</h1>
        <div class="flex flex-col border border-gray-300 rounded p-8 text-gray-900 mb-8">
        <span class="mt-1 block"><i class="material-icons text-gray-500 mr-2">drive_file_rename_outline</i><span class="font-bold text-gray-500">Description:</span> {{ eventInfo.info.description }}</span>

        <span class=""><i class="material-icons text-gray-500 mr-2">calendar_month</i><span class="font-bold text-gray-500">Date:</span> {{ eventInfo.info.dateEvent }}</span>
      </div>
        <div class="flex flex-row justify-center items-center">
          <button class="hoverButton border shadow-md border-gray-500 m-4 w-52 pl-12 pr-12 pt-2 pb-2 text-gray-500 rounded" @click="eventsRender.displayEventGoing">Going: ({{ eventInfo.going.length }})</button>
          <button class="hoverButton border shadow-md border-gray-500 m-4 w-52 pl-12 pr-12 pt-2 pb-2 text-gray-500 rounded" @click="eventsRender.displayEventNotGoing">Not going: ({{ eventInfo.notGoing.length }})</button>
        </div>
        <div class="flex flex-row justify-center items-center">
          <button @click="eventsComponent.goingToEvent(eventID)" :class="'hoverButton gradient-animation m-4 w-52 pl-12 pr-12 pt-2 pb-2 text-white font-bold rounded'">I'm going</button>
          <button @click="eventsComponent.noyGoingToEvent(eventID)" :class="'hoverButton gradient-animation m-4 w-52 pl-12 pr-12 pt-2 pb-2 text-white font-bold rounded'">I'm not going</button>
        </div>
      </div>
    </div>
    <ChatComponent></ChatComponent>
  </div>
    
    
</template>