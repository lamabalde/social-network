<script setup>
import { useRoute } from 'vue-router';
    let blockdPostMessage = 'Follow user to see posts'
    const route = useRoute().name
    if (route == 'groups') {
        blockdPostMessage = 'Join group to see posts'
    }
</script>


<template>
    <div class="items-center text-sm text-gray-900 dark:text-white">
        <p class="centerContent border border-gray-200 column bg-[white] rounded text-xl text-black">{{ blockdPostMessage }}</p>
    </div>
</template>