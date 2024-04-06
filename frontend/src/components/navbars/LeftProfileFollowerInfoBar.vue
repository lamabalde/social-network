<script setup>
import { onMounted } from 'vue';
import { followersComponent } from '../../js_modules/reactive_elements/followerComponent'
import FollowerList from './navbar_components/FollowerList.vue';
import FollowingList from './navbar_components/FollowingList.vue';
onMounted(() => {
    followersComponent.value.displayFollowBar="followers"
})
</script>

<template>
        <div class="sideBar groups">
        <div class="bg-[#ffffff]">
        <div v-if="followersComponent.displayFollowBar=='followers'">
            <FollowerList></FollowerList>
        </div>
        <div v-else-if="followersComponent.displayFollowBar=='following'">
            <FollowingList></FollowingList>
        </div>
        </div>
        </div>
</template>
