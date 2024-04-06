<script setup>
import { smallButtonstyle } from '../../../js_modules/consts';
import { routerParams } from '../../../js_modules/models/routeParamModel';
import { searchComponent } from '../../../js_modules/reactive_elements/searchComponent';
import { closeFriendsComponent } from "../../../js_modules/reactive_elements/closeFriendsComponent";
function addCloseFriend(userID) {
    closeFriendsComponent.value.addCloseFriend(userID)
}
</script>

<template>
    <ul v-if="searchComponent.closeFriendSearchResults.length" class="closeFriendsList flex bg-purple-800">
            <li v-for="user in searchComponent.closeFriendSearchResults" class="singleListElem">
                <router-link :to="new routerParams('profile', user.id)">{{ user.name }}</router-link>
                <button :class="smallButtonstyle" @click="addCloseFriend(user.id)">Add best friend</button>
            </li>
    </ul>
</template>

<style scoped>
.closeFriendsList {
    padding: 10px;
    border-radius: 10px;
    margin-top: 100px;
    margin-right: 600px;
    position: absolute;
    z-index: 0;
}
.singleListElem {
    border: solid;
    border-radius: 12px;
    padding: 5px;
}
</style>