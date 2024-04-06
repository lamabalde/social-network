<script setup>
    import { sendLogin } from '../js_modules/requests.js';
    import { ref } from 'vue';
    let error = ref()
    async function handleLogin(e) {
        error.value = await sendLogin(e)
    }
</script>

<template>
  
    <div class="grid grid-rows-[auto,1fr]">        
        <div class="flex flex-col m-16 justify-center items-center">
            <form id="loginForm" @submit.prevent="handleLogin" class="bg-slate-100 p-6 rounded-md shadow-xl w-1/4">
                <input type="text" name="username" placeholder="Username" class="w-full p-3 mb-4 border border-gray-400 rounded">
                <input type="password" name="password" placeholder="Password" class="w-full p-3 mb-4 border border-gray-400 rounded">
                <input type="submit" value="Login" class="hoverButton w-full p-3 text-white text-lg cursor-pointer gradient-animation rounded">
                <router-link to="/register" class="hoverButton flex w-full mt-10 justify-center items-center p-3 text-black text-lg cursor-pointer gradient-animation rounded">Register instead</router-link>
            </form>
        </div>
        <div v-if="error !== undefined" class="bg-red-200 border border-red-500 text-red-700 p-3 rounded m-4 text-center w-1/4 mx-auto">
            {{ error }}
        </div>
    </div>
</template>

<style>
.gradient-animation {
  background: linear-gradient(45deg, #008cff, #ec3258, #fcda00);
  background-size: 400% 400%;
  animation: gradient 6s ease infinite;
}

@keyframes gradient {
  0% {
    background-position: 0 50%;
  }
  20% {
    background-position: 20% 50%;
  }
  40% {
    background-position: 40% 50%;
  }
  60% {
    background-position: 60% 50%;
  }
  80% {
    background-position: 80% 50%;
  }
  100% {
    background-position: 0 50%;
  }
}
</style>
