<script setup>
    import { followersComponent } from '../../js_modules/reactive_elements/followerComponent';
    import { useRoute } from 'vue-router';
import ws from '../../js_modules/ws';
    const route = useRoute().params.id
    function loadFollowerStatus() {
        followersComponent.value.getFollowStatus(route)
    }
    ws.loadPageAfterWsOpen(loadFollowerStatus)
</script>

<template>
        <button v-if="followersComponent.followStatus=='following'" 
        class="gradient-animation font-base font-bold mt-10 mb-8 w-52 text-white ml-52 h-12 pl-12 pr-12 rounded" @click="followersComponent.unFollowUser(route)">
            Unfollow
        </button>
        <button v-else-if="followersComponent.followStatus=='requested'" class="gradient-animation font-base w-52 font-bold mt-10 mb-8 text-white ml-52 h-12 pl-12 pr-12 rounded">
            Requested
        </button>
        <button v-else @click="followersComponent.followUser(route)" class="gradient-animation font-base font-bold w-52 mt-10 mb-8 text-white ml-52 h-12 pl-12 pr-12 rounded">
        Follow
        </button>
</template>