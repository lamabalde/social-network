<script setup>
import { useRoute } from 'vue-router';
import { groupPageRender } from '../../../../js_modules/reactive_elements/external_view_refs/viewRender';
import { groupInfoComponent, groupMembersComponent } from '../../../../js_modules/reactive_elements/groupsComponent';
import { verifyComponent } from '../../../../js_modules/reactive_elements/verifyPageView';
import LeaveGroupButton from '../../../buttons/LeaveGroupButton.vue';
import AddGroupMemberSearch from '../../../search/AddGroupMemberSearch.vue';
const groupID = useRoute().params.id
</script>

<template>
    <div id="groupabout" v-if="groupPageRender.viewGroupAboutList" class="h-content bg-white p-4 text-black"> <!-- About the group-->
            <p class="font-bold mb-2">About us:</p>
            <p class="border-2 w-48 rounded border-purple-800 p-4">{{ groupInfoComponent.groupInfo.description }}</p>
          
            <div class="mt-8" v-if="verifyComponent.verified && groupMembersComponent.currentMemberStatus==1">
                <AddGroupMemberSearch :group-i-d="groupID"></AddGroupMemberSearch>
                <p class="font-bold mt-8">Leave group:</p>
                <LeaveGroupButton :group-i-d="groupID"></LeaveGroupButton>
            </div>
        </div>
</template>