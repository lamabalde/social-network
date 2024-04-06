<script setup>
import { eventInfo } from '../../js_modules/reactive_elements/eventsComponent';
import ListLinkElement from './navbar_components/ListLinkElement.vue';
import { eventsRender } from '../../js_modules/reactive_elements/external_view_refs/viewRender';
</script>

<template>
    <div class="sideBar groups">
        <div class="bg-[#ffffff]">
            <div v-if="eventsRender.viewEventGoing=='going'">
                <div class="flex bg-[#4a1eb3] items-center shadow-md pl-4 font-merriweather text-white h-12">
                <p class="font-medium mr-24">Going</p>
                </div>
                <ul class="p-3 overflow-y-auto h-screen"> <!-- going list -->
                    <li v-for="user in eventInfo.going" :key="user.userID" class="flex mb-2 items-center chatlist cursor-pointer">
                        <ListLinkElement :link-location="'profile'" :link-name="user.userName" :link-param="user.userID"> {{ user.userName }}</ListLinkElement>
                    </li>
                </ul> 
            </div>
        
            <div v-else-if="eventsRender.viewEventGoing=='notGoing'">
                <div class="flex bg-[#4a1eb3] items-center shadow-md pl-4 font-merriweather text-white h-12">
                <p class="font-medium mr-24">Not going</p>
                </div>
                <ul class="p-3 overflow-y-auto h-screen"> <!-- not going list -->
                    <li v-for="user in eventInfo.notGoing" :key="user.userID" class="flex mb-2 items-center chatlist cursor-pointer">
                        <ListLinkElement :link-location="'profile'" :link-name="user.userName" :link-param="user.userID"> {{ user.userName }}</ListLinkElement>
                    </li>
                </ul> 
            </div>
    </div>
    </div>
</template>