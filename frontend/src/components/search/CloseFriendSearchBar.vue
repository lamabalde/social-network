<script setup>
    import { ref, watch } from 'vue';
import { searchComponent } from '../../js_modules/reactive_elements/searchComponent';
import CloseFriendsSearchResultList from '../navbars/navbar_components/CloseFriendsSearchResultList.vue';
    const searchQuery = ref("")
    watch(() => searchQuery.value, query => {
        if (query!="") {
            searchComponent.value.searchCloseFriendsList(query)
        } else {
            searchComponent.value.closeFriendSearchResults = []
        }
    })
</script>


<template>
    <div class="flex flex-row rounded-full bg-[white] mr-8">
        <input type="text" class="rounded-full pl-4 pr-4 pt-1 pb-1 outline-none focus:outline-none text-[black]" v-model="searchQuery"
        placeholder="Add close friends">
    </div>
    
</template>

<style scoped>
.search {
  /* This bit sets up the horizontal layout */
  display:flex;
  flex-direction:row;
  /* This bit draws the box around it */
  padding: 1px;

  /* I've used padding so you can see the edges of the elements. */
}

input {
  width: 100%;
  /* Tell the input to use all the available space */
  flex-grow:2;
  /* And hide the input's outline, so the .search looks like the outline */
  border:none;
}

/* remove the input focus blue box, it will be in the wrong place. */

/* Add the focus effect to the .search so it contains the button */
input:focus-within { 
  background-color: 1px solid rgb(0, 0, 255) 
}

</style>
