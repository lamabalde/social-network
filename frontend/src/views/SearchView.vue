<script setup>
import Navbar from "../components/navbars/Navbar.vue";
import { onMounted } from 'vue';
import ws from '../js_modules/ws';
import { useRoute } from 'vue-router';
import LeftMainPageInfoBar from '../components/navbars/LeftMainPageInfoBar.vue';
import { searchComponent } from '../js_modules/reactive_elements/searchComponent';
import SearchResultList from '../components/search/SearchResultList.vue';
import ChatComponent from "../components/chat/ChatComponent.vue";
  const searchQuery = useRoute().params.id // the current group id of the route

  function loadPageData() {
    searchComponent.value.getSearchResults(searchQuery)
      //console.log("query: ",searchQuery)
  }
  onMounted(()  => {
    ws.loadPageAfterWsOpen(loadPageData)
  })
  
</script>

<template>
    <div class="page-content grid grid-cols-10 bg-[#f7f7f7]">
      <Navbar class="col-span-10"></Navbar>
      <LeftMainPageInfoBar class="col-span-2"></LeftMainPageInfoBar>
      <div id="postArea" class="col-span-6 content-center items-center">
        <SearchResultList></SearchResultList>
      </div>
  <ChatComponent/>
  </div>
</template>