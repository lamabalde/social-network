<script setup>
import { groupMembersComponent } from '../../js_modules/reactive_elements/groupsComponent';
const props = defineProps({
    groupID: String
})
function leaveGroup() {
    groupMembersComponent.value.leaveGroup(props.groupID)
}

</script>

<template>
    <button @click="leaveGroup" class="hoverButton gradient-animation mt-2 w-48 mb-4 pt-2 pb-2 text-center text-white font-bold rounded">Leave Group</button>
</template>