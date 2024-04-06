<script setup>
    import { useRoute } from 'vue-router';
import { eventsComponent } from '../../js_modules/reactive_elements/eventsComponent';
    function createEvent(ev) {
        eventsComponent.value.createGroupEvent(ev)
    }
    const groupID = useRoute().params.id
</script> 

<template>
    <div  id="addeventform" class="centerContent border border-gray-200 rounded shadow-sm justify-center bg-[#ffffff] p-16 pt-12">
        <h1 class="text-center mb-4 pb-2 pl-1 text-4xl font-medium font-merriweather">Create event</h1>
        <form class="flex flex-col" @submit.prevent="createEvent">
                <p class="mb-1">Title</p>
            <input type="text" class="m-2 rounded border-2 border-purple-500 p-2" name="title" required>
            <p class="mb-1">Description</p>
            <textarea name="description" :style="{ resize: 'none' }" class="flex m-2 rounded border-2 border-purple-500 p-2" required></textarea>
            <p class="mb-1">Day/Time</p>
            <input type="datetime-local" class="flex m-2 rounded border-2 border-purple-500 p-2" name="dateEvent" required>
            <input type="hidden" name="groupID" :value="groupID">
            <input type="submit" class="hoverButton gradient-animation mt-8 mb-4 w-52 pl-12 pr-12 pt-2 pb-2 text-white font-bold rounded" value="Create event">
        </form>
    </div>
</template>

<!-- Title       string    `json:"title"`
	Description string    `json:"description"`
	GroupID     string    `json:"groupID"`
	DateCreate  time.Time `json:"dateCreate,omitempty"`
	DateEvent   time.Time `json:"dateEvent"` -->