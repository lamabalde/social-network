<script setup>
import { eventsComponent } from '../../../../js_modules/reactive_elements/eventsComponent';
import { eventsRender } from '../../../../js_modules/reactive_elements/external_view_refs/viewRender';
import { groupMembersComponent } from '../../../../js_modules/reactive_elements/groupsComponent';
import { verifyComponent } from '../../../../js_modules/reactive_elements/verifyPageView';
import ListLinkElement from '../ListLinkElement.vue';
</script>

<template>
    <ul id="eventlist" v-if="eventsRender.viewEventsList && groupMembersComponent.currentMemberStatus==1" class="p-3 block">
        <button @click="eventsRender.viewEventForm = !eventsRender.viewEventForm" class="hoverButton gradient-animation mt-2 w-48 mb-4 pt-2 pb-2 text-center text-white font-bold rounded">Add an event</button>
            <li v-for="(event,index) in eventsComponent.eventsList" :key="event.id" class="flex mb-2 items-center chatlist cursor-pointer">
            <ListLinkElement :link-location="'event'" :link-name="event.title" :link-param="event.id"></ListLinkElement>
            </li>
    </ul> <!-- END OF GROUP EVENT LIST-->
</template>