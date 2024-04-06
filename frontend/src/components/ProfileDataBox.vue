<script setup>
  import { useRoute } from 'vue-router';
  import { currentUser, userComponent } from '../js_modules/reactive_elements/userComponent';
  import FollowButton from './buttons/FollowButton.vue';
import { followersComponent } from '../js_modules/reactive_elements/followerComponent';
import { computed, onMounted, ref} from 'vue';
import { convertDateToDateString } from '../js_modules/helpers/helpers';



const followerAmount = computed(() => followersComponent.value.followers == null ? 0 : followersComponent.value.followers.length)
const followingAmount = computed(() => followersComponent.value.following == null ? 0 : followersComponent.value.following.length)
  const route = useRoute().params.id
  //userComponent.value.userInfo.aboutMe = "I am ugly and i'm proud"
  
  function changeProfileType(profileTypeReq) {
    currentUser.value.changeProfileType(profileTypeReq)
  }
  const errorImgLoad = ref(false)


</script>


<template>
    <div v-if="userComponent.profileLoaded" class="flex flex-col bg-[#423955] h-30 w-50 rounded"> <!-- Avatar -->
        <img v-if="userComponent.userInfo.profileImg!=''" :src="'http://localhost:8000/'+userComponent.userInfo.profileImg"
        :alt="userComponent.userInfo.userName" @error="errorImgLoad=true" class="flex bg-[#ffffff] ">
        <img v-else src="/favicon.ico" class="flex bg-[#ffffff] w-40 mt-10 ml-60">
      </div> 
    <p class="mt-10 ml-8 mb-1 font-bold text-center text-4xl">
      {{ userComponent.userInfo.firstName }} 
      <span> '{{ userComponent.userInfo.userName }}' </span> 
      {{ userComponent.userInfo.lastName }} 
    </p>
    <div class="flex flex-row justify-center items-center mt-4">
        <button class="text-center text-sm" @click="followersComponent.displayFollowerBar">{{ followerAmount }} followers</button>
        <span class="mr-2 ml-2">|</span> 
        <button class="text-center text-sm" @click="followersComponent.displayFollowingBar">{{ followingAmount }} following</button>
      </div>
        <div class="border border-gray-300 rounded ml-40 mr-40 mt-10 mb-10 p-8 text-gray-900">
          <p><i class="material-icons text-gray-400 mr-1">mail</i> {{ userComponent.userInfo.email }}</p>
          <p><i class="material-icons text-gray-400 mr-1">cake</i> {{ convertDateToDateString(userComponent.userInfo.dateOfBirth)  }}</p>
        </div>
        <p class="ml-40 pb-2 text-sm">About Me: <i class="material-icons text-gray-400 text-sm">edit</i></p> <!-- Can be sent to database as userComponent.userInfo.aboutMe with Enter-->
        <textarea :placeholder="userComponent.userInfo.aboutMe || ' '" :style="{ resize: 'none' , height: 'auto' }" class="border border-gray-300 rounded ml-40 mr-40 mb-12 p-8 text-gray-900 outline-none focus:outline-none"></textarea>
        <div v-if="currentUser.userInfo.id != route">
          <FollowButton></FollowButton>
        </div>
        <div v-else>
          <p class="p-5 text-center mt-20" v-if="userComponent.userInfo.profileType==0">
            <i class="material-icons text-gray-400 text-sm">visibility</i> Public profile 
            <button class="text-gray-400" @click="changeProfileType(1)">Turn profile private</button>
          </p>
          <p class="p-5 text-center mt-20" v-else-if="userComponent.userInfo.profileType==1">
            <i class="material-icons text-gray-400 text-sm">visibility</i> Private profile 
            <button class="text-gray-400" @click="changeProfileType(0)">Turn profile public</button>
          </p>
        </div>
</template>