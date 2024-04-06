<script setup>
import { routerParams } from '../../../js_modules/models/routeParamModel';
    const props = defineProps({
        linkLocation: String,
        linkParam: String,
        linkName: String,
        avatarLink: String
    })
</script>


<template>
    <router-link :to="new routerParams(linkLocation, linkParam)" class="r-link flex cursor-pointer text-black w-full">
        <img src="/favicon.ico" class="mr-2 h-8 w-8 rounded-full"/>
        <p class="text-sm font-medium pt-1.5">{{ linkName }}</p>
    </router-link>
</template>
