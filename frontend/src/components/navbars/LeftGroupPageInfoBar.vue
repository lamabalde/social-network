<script setup>
import { groupInfoComponent } from '../../js_modules/reactive_elements/groupsComponent';
import { eventsRender, groupPageRender } from '../../js_modules/reactive_elements/external_view_refs/viewRender';
import AdminList from './navbar_components/group_info_components/AdminList.vue';
import EventList from './navbar_components/group_info_components/EventList.vue';
import GroupAbout from './navbar_components/group_info_components/GroupAbout.vue';
import GroupMemberList from './navbar_components/group_info_components/GroupMemberList.vue';


</script>

<template>
  <div class="sideBar overflow-auto members block fixed left-0 top-0 h-screen">
    <div class="bg-[#ffffff]">
        <div class="flex bg-[#4a1eb3] items-center shadow-md font-merriweather text-white h-12">
          <router-link to="/"><i class="material-icons mt-1 ml-2" title="Back">arrow_back</i></router-link><p class="font-medium pt-1 pl-3">{{ groupInfoComponent.groupInfo.name }}</p>
        </div>
        <div @click="groupPageRender.viewGroupAboutList = !groupPageRender.viewGroupAboutList" class="flex bg-[#3b188f] cursor-pointer items-center shadow-md pl-2 font-merriweather text-white h-12">
          <i class="material-icons mr-4" title="Information">info</i>      
          <p class="font-medium">Information</p>
        </div>
        <GroupAbout/>
        <div @click="eventsRender.viewEventsList = !eventsRender.viewEventsList" class="flex bg-[#2c126b] cursor-pointer items-center shadow-md pl-2 font-merriweather text-white h-12">
          <i class="material-icons mr-4" title="Event">event</i>    
          <p class="font-medium">Events</p>
        </div>
        <EventList/>
        <div @click="groupPageRender.viewAdminList = !groupPageRender.viewAdminList" class="flex bg-[#1e0c48] cursor-pointer items-center shadow-md pl-2 font-merriweather text-white h-12">
          <i class="material-icons mr-4" title="Administration">admin_panel_settings</i>     
          <p class="font-medium">Administrators</p>
        </div>
        <AdminList/>
        <div class="flex bg-[#0f0624] cursor-pointer items-center shadow-md pl-2 font-merriweather text-white h-12">
          <i class="material-icons mr-4" title="Administration">person</i>     
          <p class="font-medium">Members</p>
        </div>
        <GroupMemberList></GroupMemberList>
        <div class="mb-20">
          <!--A technical div to make overflow scrolling look better-->
        </div>
    </div>
  </div>
</template>

<style scoped>
</style>