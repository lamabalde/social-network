<script setup>
import { onMounted, watch } from 'vue';
import { routerParams } from '../../js_modules/models/routeParamModel';
import { groupListComponent } from '../../js_modules/reactive_elements/groupsComponent';
import { currentUser } from '../../js_modules/reactive_elements/userComponent';
import ListLinkElement from './navbar_components/ListLinkElement.vue';
import CreateNewGroupBox from './navbar_components/CreateNewGroupBox.vue';
watch(() => currentUser.value.userInfo.id, userID => {
    //console.log("user loaded, ", userID)
    groupListComponent.value.getGroups(userID)
})
onMounted(() => {
    if (currentUser.value.loaded) {
        groupListComponent.value.getGroups(currentUser.value.userInfo.id)
    }
})
</script>

<template>
    <div class="sideBar groups">
        <div class="bg-[#ffffff]">
        <div class="flex bg-[#4a1eb3] items-center shadow-md pl-4 font-merriweather text-white h-12">
              <p class="font-medium mr-24">My groups</p>
              <CreateNewGroupBox/>
          </div>
        <ul class="p-3 overflow-y-auto h-screen">
            <li v-for="group in groupListComponent.groups" :key="group.id" class="flex mb-2 items-center chatlist cursor-pointer">
                <ListLinkElement :link-location="'groups'" :link-name="group.name" :link-param="group.id"></ListLinkElement>
            </li>
        </ul>
    </div>
    </div>
</template>