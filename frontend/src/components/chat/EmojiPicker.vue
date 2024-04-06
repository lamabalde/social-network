<script setup>
import { onMounted, ref } from 'vue';
import { chatMessageInput } from '../../js_modules/reactive_elements/chatComponent';

function insert(emoji) {
    chatMessageInput.value += emoji.native
} 
const pickerOptions = { 
    onEmojiSelect: insert,
    theme: "light"
}
const picker = new EmojiMart.Picker(pickerOptions)

const emojiBox = ref(null)
onMounted(() => {
    emojiBox.value.appendChild(picker)
})

</script>

<template>
    <div ref="emojiBox" id="emojiBox"></div>
</template>

<style scoped>

</style>