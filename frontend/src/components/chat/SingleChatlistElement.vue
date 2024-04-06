<script setup>
import { ref } from 'vue';
import { chatComponent } from '../../js_modules/reactive_elements/chatComponent';
import { onlineUserTracker } from '../../js_modules/reactive_elements/userComponent';
    const props = defineProps({
        user: Object
    })
</script>

<template>
    <li class="mb-2 flex items-center chatlist cursor-pointer" @click="chatComponent.openCloseChatWindow(user)">
        <div class="flex flex-row items-center" v-if="user.type==0"> <!-- if chat type is private -->
            <img src="/favicon.ico" :alt="user.name" class="flex h-8 w-8 rounded-full"/>
            <p class="flex text-sm ml-1 font-base items-center content-center">
                <span class="flex items-center font-bold text-gray-600" v-if="onlineUserTracker.onlineUsers[user.id]"><i class="material-icons text-green-500 mr-1">circle</i>ONLINE:</span>
                <span class="flex items-center font-bold text-gray-600" v-else><i class="material-icons mt-1 text-red-500 mr-1">circle</i>OFFLINE:</span>
                <span class="ml-2">{{ user.name }}</span> 
                <span class="ml-1 rounded-full pr-2 pl-2 bg-gray-200">{{ user.newMessages }}</span>
            </p>
            <span v-if="user.unreadMessages>0" class="bg-red-800">{{ user.unreadMessages }}</span>
        </div>
        <div class="flex flex-row" v-else-if="user.type==1"> <!-- if chat type is group -->
            <img src="/favicon.ico" :alt="user.name" class="flex h-8 w-8 rounded-full"/>
            <p class="flex text-sm ml-1 font-base items-center content-center">
                <span class="flex font-bold text-gray-600"><i class="material-icons mr-1 text-gray-500">circle</i>GROUP:</span>
                <span class="ml-2">{{ user.name }}</span>
                <span v-if="user.newMessages>0" class="ml-1 rounded-full pr-2 pl-2 bg-red-800">{{ user.newMessages }}</span>
            </p>
            <span v-if="user.unreadMessages>0" class="bg-red-800">{{ user.unreadMessages }}</span>
        </div>
    </li>
</template>