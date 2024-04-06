<script setup>
import { followersComponent } from '../../../js_modules/reactive_elements/followerComponent';
import ListLinkElement from './ListLinkElement.vue';
</script>

<template>
    <div class="flex bg-[#4a1eb3] items-center shadow-md pl-4 font-merriweather text-white h-12">
        <p class="font-medium">Following</p>
    </div>
    <ul class="p-3 overflow-y-auto h-screen">
        <li v-for="follower in followersComponent.following" class="flex mb-2 items-center chatlist cursor-pointer">
            <ListLinkElement :link-location="'profile'" :link-name="follower.userName" :link-param="follower.id"></ListLinkElement>
        </li>
    </ul>
</template>