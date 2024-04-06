<script setup>
import { ref } from 'vue';
import { groupListComponent } from '../../../js_modules/reactive_elements/groupsComponent';
    
    
let showNewGroupBox = ref(false)
function openNewGroupBox() {
    showNewGroupBox.value = !showNewGroupBox.value
}

function groupcreated() {
  var eventList = document.getElementById("confirmdone");
  if (eventList.style.display === 'block' || eventList.style.display === '') {
    eventList.style.display = 'none';
  } else {
    eventList.style.display = 'block';
  }
}
</script>

<template>
    <button class="pt-2" title="Add new group" @click="openNewGroupBox"><i class="material-icons">group_add</i></button>
    <div v-if="showNewGroupBox" :style="{ width : '238px', fontFamily : 'inter' , top: '115px' }" class="fixed pl-4 pt-4 text-black left-0 bg-white h-screen ml-0">
        <form @submit.prevent="groupListComponent.createGroup">
            <p class="mb-4 text-center font-bold font-merriweather">Create a new group</p>
            <p class="mb-1">Title</p>
        <input type="text" :style="{ width: '200px' }" class="pl-2 pr-2 mb-4 rounded border-2 border-purple-500" name="title" required>
        <p class="mb-1">Description</p>
        <textarea name="description" :style="{ resize: 'none' , width: '200px' }" class="p-2 h-20 rounded border-2 border-purple-500" required></textarea>
        <input type="submit" class="hoverButton gradient-animation mt-8 mb-4 pl-12 pr-12 pt-2 pb-2 text-white font-bold rounded" @click="groupcreated()" value="Create group">
        <p id="confirmdone" :style="{ display: 'none' }" class="mb-1 ml-12 text-[#4cc218]">Group created!</p>
    </form>
    </div>
    
</template>
