<script setup>
import { postsComponent } from '../../js_modules/reactive_elements/postComponent';
import Post from './PostBox.vue';
</script>

<template>
    <div class="posts items-center text-sm text-gray-900 dark:text-white">
        <Post v-for="post in postsComponent.posts" :key="post.id" :title="post.theme" 
        :body="post.content.text" :userName="post.content.userName" :postID="post.id" :userID="post.content.userID" 
        :dateCreate="post.content.dateCreate" :likes="post.content.likes[1]" :post="post"/>
    </div>
    
</template>