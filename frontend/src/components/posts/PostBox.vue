<script setup>
import { convertDateToLocale } from '../../js_modules/helpers/helpers';
import { routerParams } from '../../js_modules/models/routeParamModel';
import { postsComponent } from '../../js_modules/reactive_elements/postComponent';
import { currentUser } from '../../js_modules/reactive_elements/userComponent';
import PostComments from '../comments/PostComments.vue';
import PostLikes from '../reactions/PostLikes.vue';
    const props = defineProps({
        post: Object,
        title: String,
        body: String,
        userName: String,
        postID: String,
        userID: String,
        dateCreate: String,
        likes: Number,
        dislikes: Number
    })
    
</script>

<template>
        <div class="centerContent border border-gray-200 column bg-[white] rounded text-xl text-black"> <!-- single postbox style -->
            <button v-if="userID === currentUser.userInfo.id" @click="postsComponent.deletePostReq(postID)" :style="{ right: '25%' }" class="flex absolute text-center p-3 mt-2 text-lg rounded cursor-pointer text-gray-300 hover:text-gray-500">
            Delete post <i class="material-icons">delete_forever</i>
            </button>
            <div class="flex flex-col ml-8 mt-6 mr-8 mb-8">
                <div class="flex flex-row mb-8"> <!-- User info -->
                    <img src="/favicon.ico" class="h-14 w-14 m-1 rounded-full border-2 border-purple-500"/>
                    <div class="flex flex-col ml-1 mt-1">
                        <router-link :to="new routerParams('profile', userID)" class="font-bold text-l">{{ userName }}</router-link>
                        <span class="text-sm"> {{ convertDateToLocale(dateCreate)   }} <!-- 13.12.2023 14:49 --></span>
                    </div>
                </div>
            <div class=""> <!-- Post content -->
                <h2 class="font-bold mb-8 text-3xl font-merriweather">{{ title }}</h2>
                <p class="">{{ body }}</p>
            </div>
            </div>
            <img v-if="post.image" :src="'http://localhost:8000/'+post.image">
            <PostLikes :likes="likes" :post-i-d="postID"></PostLikes>
            <PostComments :post-i-d="postID"></PostComments>
        </div>
</template>

<style scoped>
</style>