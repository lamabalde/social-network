<script setup>
import { useRoute } from 'vue-router';
import { groupMembersComponent } from '../../js_modules/reactive_elements/groupsComponent';
const groupID = useRoute().params.id

const props = defineProps({
    userID: String,
    userName: String
})
function inviteToGroup() {
    groupMembersComponent.value.inviteUserToGroup(groupID, props.userID)
}

</script>

<template>
    <button @click="inviteToGroup" class="gradient-animation ml-2 pl-2 pr-2 text-white rounded">Invite To Group</button>
</template>