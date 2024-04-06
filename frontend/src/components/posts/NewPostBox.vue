<script setup>
    import { useRoute } from 'vue-router';
    import { postsComponent } from '../../js_modules/reactive_elements/postComponent';
    const props = defineProps(
        {
            newPostText: String,
        }
    )
    const route = useRoute().params.id
    function handleSubmitPost(e) {
        postsComponent.value.submitNewPost(e)
        e.target.reset()
        
    }
    
</script>

<template>
    <div class="centerContent flex flex-col border border-gray-200 rounded shadow-sm justify-center bg-[#ffffff] p-16 pt-12">
        <h1 class="text-center mb-4 pb-2 pl-1 text-4xl font-medium font-merriweather">{{ newPostText }}</h1>
        <form id="submitPost" @submit.prevent="handleSubmitPost" class="flex flex-col">
            <input type="text" name="Title" placeholder="Title" class="m-2 rounded border-2 border-purple-500 p-2" required />
            <textarea name="body" placeholder="Content" class="m-2 rounded border-2 border-purple-500 p-2" rows="4" required></textarea>
            <div v-if="useRoute().name!='groups'">
                <label for="postPrivacy">Select post privacy:</label>
                    <select name="postType" id="postPrivacy" >
                        <option :value="0" selected>Public</option> <!-- post is public by default -->
                        <option :value="1">Private</option>
                        <option :value="2">Friends</option>
                    </select>
            </div>
            
            <label for="postImage">Select an image:</label>
            <input type="file" name="image" placeholder="Select an Image" id="postImage"
                class="w-full p-3 mb-4 border-2 border-purple-500 rounded" accept="image/*">
            <input v-if="useRoute().name=='groups'" type="hidden" name="groupID" :value="route">
            <input type="submit" value="Post" :class="'flex hoverButton items-center justify-center text-center font-bold w-52 ml-40 mr-40 mt-8 gradient-animation rounded cursor-pointer h-10 text-white'" />
        </form>
    </div>
</template>

<style scoped>
label {
    margin: 5px;
}
</style>