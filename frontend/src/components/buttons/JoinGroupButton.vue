<script setup>
import { groupMembersComponent } from '../../js_modules/reactive_elements/groupsComponent';
const props = defineProps({
    groupID: String
})
function joinGroup() {
    groupMembersComponent.value.joinGroup(props.groupID)
}

</script>

<template>
    <button @click="joinGroup" class="centerContent gradient-animation text-white ml-60 h-12 pl-12 pr-12 rounded-full">Request to join group</button>
</template>