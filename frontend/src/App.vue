<script setup>
import { watch } from 'vue';
import { notificationsComponent } from './js_modules/reactive_elements/notificationComponent';
import { currentUser } from './js_modules/reactive_elements/userComponent';
import ChatComponent from './components/chat/ChatComponent.vue';
import Chat from './components/chat/Chat.vue';
import { chatComponent } from './js_modules/reactive_elements/chatComponent';
watch(() => currentUser.value.userInfo.id, id => {
  getNotifications(id)
})

function getNotifications(userID) {
  notificationsComponent.value.getNotifications(userID)
}
</script>


<template>
  <router-view :key="$route.fullPath"/>
  <Chat v-if="chatComponent.isChatOpen"></Chat>
</template>