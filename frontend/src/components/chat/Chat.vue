<script setup>
  import { chatComponent, chatMessageInput } from '../../js_modules/reactive_elements/chatComponent';
  import { currentUser } from '../../js_modules/reactive_elements/userComponent';
import { onMounted,onUnmounted,ref, watch } from 'vue';
import {atTop} from '../../js_modules/helpers/scrollListener'
import { convertDateToDateString, convertDateToLocale } from '../../js_modules/helpers/helpers';
import EmojiPicker from './EmojiPicker.vue';
  let oldScrollHeight = 0
  const chatBox = ref(HTMLElement | null)
  const recipient = ref(chatComponent.value.currentOpenChat)
  function sendMessage() {
    chatComponent.value.sendMessage(chatComponent.value.currentOpenChat.id, chatMessageInput.value)
    chatMessageInput.value = ""
  }

  // get the user id from the chatlist here
  watch(() => atTop.value, scrolledTop => {
    if (scrolledTop) {
      oldScrollHeight = chatBox.value.scrollHeight
      chatComponent.value.offset+=10
      chatComponent.value.getChatMessages(recipient.value.id)
      watch(() => chatComponent.value.messages, updated => { // have it retain the current position of the message
        if (chatBox.value!= null) {
          chatBox.value.scrollTo({
          top: chatBox.value.scrollHeight - oldScrollHeight
        })
        }
        
      })
    } 
  })
  
  watch(() => chatComponent.value.messages.length, messages =>setTimeout(() => { // TODO: fix error when closing chat
    chatBox.value.scrollTo({ // scroll to bottom of chat when messages have loaded
      top: chatBox.value.scrollHeight,
      left:0,
      behavior: "smooth"
    })
  },10))
  
  onMounted(() =>{
    watch(() => chatComponent.value.messagesLoaded, loaded => { // TODO: fix error when closing chat
      chatBox.value.scrollTo({ // scroll to bottom of chat when messages have loaded
        top: chatBox.value.scrollHeight,
        left:0,
        behavior: "smooth"
      })
    })
  })
  function handleScroll(e) {
    if (e.target.scrollTop<=0 && !atTop.value) {
      atTop.value=true
    } else {
      atTop.value=false
    }
  }
  const emojiBoxOpen = ref(false)
  function openEmojiBox() {
    emojiBoxOpen.value = !emojiBoxOpen.value
  }
</script>

<template>
  <div class="flex flex-row"> <!-- alignment of emoji box and chat box -->
    <div class="fixed bottom-3 right-3 flex flex-col border rounded" :style="{height: '350px', width: '380px'}" ref="chatDiv">
    <!-- Top Area -->
    <div :style="{ borderTopLeftRadius: '4px', borderTopRightRadius: '4px'}" class="flex items-center justify-between bg-[#4a1eb3] px-4 py-2 text-white shadow-md font-merriweather h-12">
      <span class="text-lg font-medium">{{ recipient.name }} </span>
      <button class="text-xl font-semibold" @click="chatComponent.openCloseChatWindow(recipient)"><i class="material-icons mt-2">close</i></button>
    </div>

    <!-- Center Area for messages -->
    <div class="flex-1 overflow-y-scroll p-5 bg-gray-200" id="chatBox" ref="chatBox" @scroll="handleScroll"> <!-- TODO need to add a scroll lietener here for message offset -->
      <div v-for="(message, index) in chatComponent.messages" :key="index" class="flex flex-col">
        <div v-if="message.userID == currentUser.userInfo.id" class="border mb-1 p-2 rounded bg-purple-300">
          <strong >{{ message.userName }}:</strong> {{ message.content }} {{ convertDateToLocale(message.dateCreate)  }}
        </div>
        <div v-else class="flex-end border mb-1 p-2 rounded bg-white">
          <strong >{{ message.userName }}:</strong> {{ message.content }} {{ convertDateToLocale(message.dateCreate)  }}
        </div>
      </div>
    </div>

    <!-- Input Box and Send Button -->
    <div :style="{ borderBottomLeftRadius: '4px', borderBottomRightRadius: '4px'}" class="bg-gray-100 flex items-center justify-between p-4"  @keyup.enter="sendMessage">
     
      <input v-model="chatMessageInput" class="flex-1 rounded border p-2" placeholder="Type your message..." />
      <span class="material-icons cursor-pointer" @click="openEmojiBox">mood</span>
      <button class="hoverButton gradient-animation rounded ml-2 bg-purple-800 px-4 py-2 text-white"  @click="sendMessage">Send</button>
    </div>
  </div>
  <EmojiPicker class="fixed" :style="{ right : '36%', bottom : '1%'}" v-if="emojiBoxOpen"></EmojiPicker> <!-- end-96 -->
  </div>
  
</template>