<script setup>
  import { watch } from 'vue';
  import { chatComponent } from '../../js_modules/reactive_elements/chatComponent';
  import { currentUser, onlineUserTracker } from '../../js_modules/reactive_elements/userComponent';
import SingleChatlistElement from './SingleChatlistElement.vue';

  
  watch(() => currentUser.value.userInfo.id, UpdatedUserID => { //getChatList when userID loads
    chatComponent.value.getChatList(UpdatedUserID)
  })
</script>

<template>
  <div :style="{ borderLeft: '1px solid #dfdfdf'}" class="fixed right-0 mt-16 w-60">
    <div class="bg-[#ffffff]">
      <div class="flex bg-[#4a1eb3] items-center shadow-md pl-4 font-merriweather text-white h-12">
        <p class="font-medium">Chats</p>
      </div>
      <ul class="p-3 overflow-y-auto h-screen">
        <SingleChatlistElement v-for="user in chatComponent.chatList" :key="user.id"  :user="user"></SingleChatlistElement>

      </ul>
    </div>
  </div>
</template>

<style scoped>
li:last-child {
    margin-bottom: 7rem;
}
::-webkit-scrollbar {
    width: 10px;
}
 
::-webkit-scrollbar-track {
    background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 5px;
}

::-webkit-scrollbar-thumb:hover {
    background: #555;
}
</style>