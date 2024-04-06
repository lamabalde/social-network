<script setup>
  import Navbar from "../components/navbars/Navbar.vue";
  import LeftProfileFollowerInfoBar from '../components/navbars/LeftProfileFollowerInfoBar.vue';
  import ws from '../js_modules/ws';
  import { onMounted, ref, watch } from 'vue';
  import { useRoute } from 'vue-router';
  import { postsComponent } from '../js_modules/reactive_elements/postComponent';
  import { followersComponent } from '../js_modules/reactive_elements/followerComponent'
  import { userComponent } from '../js_modules/reactive_elements/userComponent';
  import PostList from '../components/posts/PostList.vue';
  import ProfileDataBox from "../components/ProfileDataBox.vue";
  import BlockedPostList from "../components/posts/BlockedPostList.vue";
 import { atBottom, scrollAtBottomListener } from '../js_modules/helpers/scrollListener';
  import {verifyComponent} from "../js_modules/reactive_elements/verifyPageView"
import ChatComponent from "../components/chat/ChatComponent.vue";
  userComponent.value.userInfo.profileImg = "" // to prevent wrong images loading
  const route = useRoute().params.id
  const offset = ref(0) // when scrolling down to post limit, it will add 10 to offset and launch getAllPosts again
  function getNewPosts() { // have this activate when at the bottom of posts
      postsComponent.value.getUserPosts(route, offset.value)
  }
  const canSeePosts = verifyComponent.value.verifyFollowView()
  function loadPage() {
    userComponent.value.getUserInfo(route)
    followersComponent.value.getFollowers(route)
    followersComponent.value.getFollowing(route)
    getNewPosts()
    scrollAtBottomListener()
  }
  onMounted(()  => {
    ws.loadPageAfterWsOpen(loadPage)
  })
  watch(() => atBottom.value, scrolledBottom => {
    if (scrolledBottom) {
      offset.value+=10
      getNewPosts()
    }
  })
</script>

<template>
  <div class="main bg-[#f7f7f7] flex flex-col">
  <Navbar></Navbar>
  <LeftProfileFollowerInfoBar></LeftProfileFollowerInfoBar>
    <div id="pageInfo" class="centerContent flex flex-col bg-[#ffffff] border border-gray-200 rounded mt-16 w-20 h-100">

      <ProfileDataBox></ProfileDataBox>
        
    </div>
    <PostList v-if="canSeePosts"></PostList>
    <BlockedPostList v-else></BlockedPostList>
    <ChatComponent></ChatComponent>
    </div>
</template>
