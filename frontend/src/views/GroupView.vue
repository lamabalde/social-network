<script setup>
import NewPostBox from '../components/posts/NewPostBox.vue';
import Navbar from "../components/navbars/Navbar.vue";
import { onMounted, ref, watch } from 'vue';
import LeftGroupPageInfoBar from '../components/navbars/LeftGroupPageInfoBar.vue';
import ws from '../js_modules/ws';
import { useRoute } from 'vue-router';
import { postsComponent } from '../js_modules/reactive_elements/postComponent';
import { groupMembersComponent } from '../js_modules/reactive_elements/groupsComponent';
import PostList from '../components/posts/PostList.vue';
import {verifyComponent} from '../js_modules/reactive_elements/verifyPageView'
import BlockedPostList from '../components/posts/BlockedPostList.vue';
import { groupInfoComponent } from '../js_modules/reactive_elements/groupsComponent';
import JoinGroupButton from '../components/buttons/JoinGroupButton.vue';
import LeaveGroupButton from '../components/buttons/LeaveGroupButton.vue';
import { atBottom, scrollAtBottomListener } from '../js_modules/helpers/scrollListener';
import RequestedButton from '../components/buttons/RequestedButton.vue';
import AddEvent from '../components/events/AddEvent.vue';
import { eventsRender } from '../js_modules/reactive_elements/external_view_refs/viewRender';
import { eventsComponent } from '../js_modules/reactive_elements/eventsComponent';
import ChatComponent from "../components/chat/ChatComponent.vue";
const offset = ref(0) // when scrolling down to post limit, it will add 10 to offset and launch getAllPosts again
const groupID = useRoute().params.id // the current group id of the route
  const newPostText="Add post to group"
  verifyComponent.value.finishedVerification = false

  function initialVerify() {
    verifyComponent.value.verifyGroupView(groupID)
    groupInfoComponent.value.getGroupInfo(groupID)
    groupMembersComponent.value.getMembers(groupID)
  }
  onMounted(() => {
    ws.loadPageAfterWsOpen(initialVerify)
    
  })
  
  function getNewPosts() { // have this activate when at the bottom of posts
    postsComponent.value.getGroupPosts(groupID, offset.value)
  }
  function loadPageData() {
      getNewPosts()
      scrollAtBottomListener()
      eventsComponent.value.getGroupEvents(groupID)
  }
  
  watch(()=> verifyComponent.value.finishedVerification, verify => {
    if (verifyComponent.value.verified) {
      ws.loadPageAfterWsOpen(loadPageData)
    }
  })
  watch(() => atBottom.value, scrolledBottom => {
    if (scrolledBottom) {
      offset.value+=10
      getNewPosts()
    }
  })
</script>

<template>
  <div id="mainPage" class="content-center flex flex-col bg-[#f7f7f7] h-screen">

    <Navbar></Navbar>
    
    <div class="flex flex-row bg-[#f7f7f7] mt-8">

      <LeftGroupPageInfoBar></LeftGroupPageInfoBar>
      <div id="allPostData" class="content-center flex flex-col">
        <AddEvent v-if="eventsRender.viewEventForm"></AddEvent>
        <div v-if="verifyComponent.verified && groupMembersComponent.currentMemberStatus==1" >
          <LeaveGroupButton :group-i-d="groupID"></LeaveGroupButton>
          <NewPostBox :new-post-text="newPostText"/>
          <PostList></PostList>
        </div>
        <div v-else>
          <BlockedPostList></BlockedPostList>
          
          <RequestedButton v-if="groupMembersComponent.currentMemberStatus==0"></RequestedButton>
          <JoinGroupButton v-else :group-i-d="groupID" ></JoinGroupButton>
        </div>
      </div>
      
    </div>
      
    
    <ChatComponent/>
  </div>
</template>