<script setup>
import { routerParams } from '../../js_modules/models/routeParamModel';
import { searchComponent } from '../../js_modules/reactive_elements/searchComponent';

</script>

<template>
    <ul class="mt-20 ">
        <li v-for="result in searchComponent.searchResults" :key="result.id" class="m-5">
                <router-link :to="new routerParams(result.entityType , result.id)" class="centerContent flex-auto column shadow-xl bg-gradient-to-r from-cyan-500 to-purple-500 rounded text-xl">
                    {{ result.name }}
                </router-link>
            </li>
    </ul>
    
</template>