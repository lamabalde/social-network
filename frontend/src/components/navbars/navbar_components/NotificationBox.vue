<script setup>
import { ref } from 'vue';
import { notificationsComponent } from '../../../js_modules/reactive_elements/notificationComponent';
import SingleNotification from './SingleNotification.vue';
const showNotifications = ref({
    visible: false,
    show() {
        this.visible = !this.visible
    }
})
</script>


<template>
    <button class="mr-4" title="Notifications" @click="showNotifications.show"><i class="material-icons">notifications</i></button> <!-- TODO: will need to get current user notifications in a reactive object -->
    <ul v-if="showNotifications.visible" class="absolute top-16 text-black p-4 rounded border-2 bg-[#f7f7f7] border-purple-800">
        <SingleNotification v-if="notificationsComponent.notifications.length" 
        v-for="notification in notificationsComponent.notifications" :notification="notification" class="border border-gray-200 bg-white m-2 p-2 rounded flex flex-col h-content">
        </SingleNotification>
        <span v-else class="flex h-content">No new notifications</span>
    </ul>
</template>
