<script setup>
import { postsComponent } from '../../js_modules/reactive_elements/postComponent';

const props = defineProps({
    likes: Number,
    postID: String
})
function likePost() {
    postsComponent.value.likePost(props.postID)
}
</script>

<template>
    <div class="flex flex-row ml-8">
        <i class="material-icons text-red-500 mr-2 cursor-pointer" @click="likePost">favorite_border</i> <!-- Change material icon to "favorite" after click -->
        <span>{{ likes }} likes</span>
    </div>
</template>