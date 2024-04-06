<script setup>
    import { onMounted, ref, toRefs } from 'vue';
    import { commentsComponent } from '../../js_modules/reactive_elements/commentsComponent';
import { convertDateToLocale } from '../../js_modules/helpers/helpers';
    const props = defineProps({
        postID: String
    })
    let commentsVisible= ref(false) // this needs to be inside this components because we don't want a single source of truth for it
    onMounted(() => {
        commentsComponent.value.getAllComments(props.postID)
    })

    const showComments = () => {
        commentsVisible.value = !commentsVisible.value
    }
    function submitComment(ev) {
        
        commentsComponent.value.submitComment(ev)
        ev.target.reset()
    }
</script>

<template>
    <button class="ml-8 mt-10 hover:text-purple-800" @click="showComments">Comments ({{ commentsComponent.amountOfComments(postID) }}) </button>
    <div v-if="commentsVisible">
        <div v-for="comment in commentsComponent.comments[postID]" class="flex flex-row m-8 border rounded p-4 ">
            <img src="/favicon.ico" class="h-8 w-8 m-1 rounded-full border-2 border-purple-500"/>
            <div class="flex flex-col ml-2 mt-1">
                
                <span class="mb-2 font-bold text-gray-600">{{ comment.content.userName }}</span>
                <span class="mb-2">{{ comment.content.text }}</span>
                <img v-if="comment.content.image" :src="'http://localhost:8000/'+comment.content.image">
                <span class="mb-2"> {{ convertDateToLocale(comment.content.dateCreate)  }} </span>
            </div>
        </div>
    </div>
    
    <div class="flex flex-row m-8 border-2 border-purple-500 rounded">
        
    <form @submit.prevent="submitComment" class="flex flex-col">
        
        <textarea name="content" placeholder="Write a comment..."
         :style="{ width: '26.8vw' }" class="m-5 h-8 text-black outline-none focus:outline-none" required></textarea>
        <!-- <input type="text" name="content" placeholder="Write a comment..." 
        :style="{ width: '27.8vw' }" class="m-5 text-black outline-none focus:outline-none" required> -->
        <input type="hidden" name="postID" :value="postID" >
        <input type="file" name="image" accept="image/*">
        <input type="submit" value="Comment" :style="{ right: '28%' }" class="flex absolute m-3 pl-10 text-base text-center hoverButton font-bold w-40 gradient-animation rounded cursor-pointer h-10 text-white">
        
    </form>
    </div>
</template>