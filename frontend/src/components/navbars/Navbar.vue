<script setup>
import NotificationBox from './navbar_components/NotificationBox.vue';
import LogoutButton from './navbar_components/LogoutButton.vue';
import { routerParams } from '../../js_modules/models/routeParamModel';
import { currentUser } from '../../js_modules/reactive_elements/userComponent';
import SearchBar from '../search/SearchBar.vue';
import CloseFriendSearchBar from '../search/CloseFriendSearchBar.vue';
import CloseFriendsSearchResultList from './navbar_components/CloseFriendsSearchResultList.vue';
</script>

<template>
  <header class="navbar flex items-center w-full fixed shadow-md top-0 bg-purple-800 text-white text-m h-16">
    <router-link to="/">
        <img src="/favicon.ico" title="Go to the main page" alt="home" class="w-16 ml-4">
      </router-link>
      
    <div class="navbar-items flex justify-end items-center pl-4 pr-6 w-full">
      <CloseFriendSearchBar>
      </CloseFriendSearchBar>
      <CloseFriendsSearchResultList></CloseFriendsSearchResultList>
      <SearchBar></SearchBar>
      <!-- the v-if is so that the router doesn't crash when there is an empty string -->
      <div class="flex justify-between mr-8">
        
      <router-link class="flex" title="Profile" v-if="currentUser.userInfo.id!=''" :to="new routerParams('profile', currentUser.userInfo.id)">
        {{ currentUser.userInfo.userName }} <i class="material-icons mr-4 ml-4">account_circle</i></router-link>
      
      <NotificationBox/>
    </div>
      <LogoutButton/>
    </div>
  </header>
</template>

<style scoped>
.navbar {
  z-index: 3;
}
</style>