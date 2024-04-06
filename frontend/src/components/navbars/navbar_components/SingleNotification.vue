<script setup>
import { RouterLink } from 'vue-router';
import { followersComponent } from '../../../js_modules/reactive_elements/followerComponent';
import { routerParams } from '../../../js_modules/models/routeParamModel';
    import { groupMembersComponent } from '../../../js_modules/reactive_elements/groupsComponent';
import { convertDateToLocale } from '../../../js_modules/helpers/helpers';
    const props = defineProps({
        notification: Object
    })

    const notificationID = props.notification.id
    function acceptFollow() {
        followersComponent.value.acceptFollowRequest({id: notificationID, userID: props.notification.fromUserId.String} )
    }
    function declineFollow() {
        followersComponent.value.declineFollowRequest({id: notificationID, userID: props.notification.fromUserId.String})
    }
    function acceptGroupInvite() {
        groupMembersComponent.value.acceptGroupInvite(notificationID)
    }
    function declineGroupInvite() {
        groupMembersComponent.value.declineGroupInvite(notificationID)
    }
    function acceptGroupRequest() {
        groupMembersComponent.value.acceptGroupJoinRequest(notificationID)
    }
    function declineGroupRequest() {
        groupMembersComponent.value.declineGroupJoinRequest(notificationID)
    }
</script>



<template>
    <div v-if="notification.type == 'follow request'" class="flex flex-col">
        <RouterLink :to="new routerParams('profile', notification.fromUserId.String)">
            {{  notification.body }}
            
        </RouterLink>
        <span>{{ convertDateToLocale(notification.dateCreate) }}</span> 
        <div class="flex-row">
            
            <button class="text-center pr-2 pl-2 mr-2 font-bold gradient-animation rounded cursor-pointer text-white" @click="acceptFollow">Accept</button>
            <button class="rounded pr-2 pl-2 border shadow-sm border-gray-500" @click="declineFollow">Decline</button>
        </div>
    </div>
    <div v-else-if="notification.type == 'invite to group'" class="flex flex-col">
        <RouterLink :to="new routerParams('profile', notification.fromUserId.String)">{{ notification.body }}</RouterLink>
        <span>{{ convertDateToLocale(notification.dateCreate) }}</span> 
        <div class="flex-row">
            <button class="text-center pr-2 pl-2 mr-2 font-bold gradient-animation rounded cursor-pointer text-white" @click="acceptGroupInvite">Accept</button>
            <button class="rounded pr-2 pl-2 border shadow-sm border-gray-500" @click="declineGroupInvite">Decline</button>
        </div>
    </div>
    <div v-else-if="notification.type == 'join group request'" class="flex flex-col">
        <RouterLink :to="new routerParams('profile', notification.fromUserId.String)">{{ notification.body }}</RouterLink>
        <span>{{ convertDateToLocale(notification.dateCreate) }}</span> 
        <div class="flex-row">
            <button class="text-center pr-2 pl-2 mr-2 font-bold gradient-animation rounded cursor-pointer text-white" @click="acceptGroupRequest">Accept</button>
            <button class="rounded pr-2 pl-2 border shadow-sm border-gray-500" @click="declineGroupRequest">Decline</button>
        </div>
    </div>
    <div v-else>
        <p>{{ notification.body }}</p>
        <span>{{ convertDateToLocale(notification.dateCreate) }}</span> 
    </div>
    
</template>

<style scoped>
</style>