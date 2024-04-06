<script setup>
import { groupInfoComponent, groupMembersComponent } from '../../../../js_modules/reactive_elements/groupsComponent';
import ListLinkElement from '../ListLinkElement.vue';
</script>

<template>
    <ul id="memberlist" class="p-3 h-content">
          <li v-for="(member,index) in groupMembersComponent.members" :key="member.id" class="flex mb-2 items-center chatlist cursor-pointer">
            <ListLinkElement :link-location="'profile'" :link-name="member.userName" :link-param="member.id" v-if="member.id!=groupInfoComponent.groupInfo.creator_id"></ListLinkElement>
          </li>
        </ul>
</template>