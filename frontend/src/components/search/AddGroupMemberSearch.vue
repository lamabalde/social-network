<script setup>
import { ref, watch } from 'vue';
import { searchComponent } from '../../js_modules/reactive_elements/searchComponent';
import InviteToGroupButton from '../buttons/InviteToGroupButton.vue';
import { currentUser } from '../../js_modules/reactive_elements/userComponent';
    const props = defineProps({
        groupID: String
    })


    searchComponent.value.searchResults=[]
    const searchQuery = ref("")
    watch(() => searchQuery.value, query => { // make ws request to fetch users on every letter
        if (query == "") {
            searchComponent.value.searchResults=[]
        } else {
            searchComponent.value.getUserSearchResultsForGroup(query, props.groupID)
        }
    })
</script>

<template>
    <div class="flex flex-col rounded bg-[#ffffff]">
        <p class="font-bold mb-2">Add user to group:</p>
        <input type="text" v-model="searchQuery" class="border-2 rounded w-48 border-purple-800" placeholder="Search by username">
        <ul class="bg-[#ffffff] w-48">
            <li v-for="user in searchComponent.searchResults" class="pl-2 pr-2 pt-1 pb-1">
                <div v-if="user.id != currentUser.userInfo.id">
                    {{ user.userName }} <InviteToGroupButton :user-i-d="user.id" :user-name="user.name"></InviteToGroupButton>
                </div>
            </li>
        </ul> 
    </div>
</template>