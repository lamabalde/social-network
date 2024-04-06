<script setup>
  import { onMounted, ref, watch } from 'vue';
  import ws from "../js_modules/ws";
  import NewPostBox from '../components/posts/NewPostBox.vue';
  import Navbar from "../components/navbars/Navbar.vue";
  import LeftInfoBar from '../components/navbars/LeftMainPageInfoBar.vue';
  import { useRoute } from 'vue-router';
  import { postsComponent } from '../js_modules/reactive_elements/postComponent';
  import PostList from '../components/posts/PostList.vue';
  import { atBottom, scrollAtBottomListener } from '../js_modules/helpers/scrollListener';
  import ChatComponent from "../components/chat/ChatComponent.vue";
  
  const route = useRoute().params.id // the current id of the route
  const offset = ref(0) // when scrolling down to post limit, it will add 10 to offset and launch getAllPosts again
  function getNewPosts() { // have this activate when at the bottom of posts
    postsComponent.value.getMainPosts(offset.value)
  }
  function loadPageInfo() {
    getNewPosts()
    scrollAtBottomListener()
  }
  
  watch(() => atBottom.value, scrolledBottom => {
    if (scrolledBottom) {
      offset.value+=10
      getNewPosts()
    }
  })
  
  onMounted(() => {
    ws.loadPageAfterWsOpen(loadPageInfo)
  })
  

  function onScroll ({ target: { scrollTop, clientHeight, scrollHeight }}) {
    console.log(target)
      if (scrollTop + clientHeight >= scrollHeight) {
        console.log("jeees")
      }
    }
  const newPostText="Add new post"

  
</script>

<template>
  <div class="page-content bg-[#f7f7f7] grid grid-cols-10" @scroll="onScroll">
      <Navbar class="col-span-10"></Navbar>
      <LeftInfoBar class="col-span-2"></LeftInfoBar>
      <div id="postArea" class="col-span-6 content-center items-center mt-16" >
        <NewPostBox :new-post-text="newPostText"/>
        <PostList></PostList>
      </div>
      <ChatComponent/>
  </div>
</template>