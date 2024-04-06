<script setup>
import { ref } from 'vue';
import { sendRegister } from '../js_modules/requests.js';
    let error = ref()
    async function handleRegister(e) {
        
        error.value = await sendRegister(e)
    }
</script>

<template>
    
    <body class="overflow-hidden">
    <div class="grid grid-rows-[auto,1fr] m-8">
        <div class="flex flex-col items-center justify-center">
            <form id="registerForm" @submit.prevent="handleRegister" class="bg-slate-100 p-6 rounded-md shadow-xl w-1/4 mx-auto ">
                <input type="email" name="email" placeholder="Email*" class="w-full p-3 mb-4 border border-gray-400 rounded" required>
                <input type="text" name="firstName" placeholder="First Name*" class="w-full p-3 mb-4 border border-gray-400 rounded" required>
                <input type="text" name="lastName" placeholder="Last Name*" class="w-full p-3 mb-4 border border-gray-400 rounded" required>
                <p>Date of birth:</p>
                <input type="date" name="dateOfBirth" placeholder="Date of Birth*" class="w-full p-3 mb-4 border border-gray-400 rounded text-gray-500" required>
                <input type="text" name="aboutMe" placeholder="About Me" class="w-full p-3 mb-4 border border-gray-400 rounded">
                <input type="text" name="userName" placeholder="Username" class="w-full p-3 mb-4 border border-gray-400 rounded">
                <input type="password" name="password" placeholder="Password*" class="w-full p-3 mb-4 border border-gray-400 rounded">
                <p>Select a profile image:</p>
                <input type="file" name="image" placeholder="Select an Image" 
                class="w-full p-3 mb-4 border border-gray-400 rounded" accept="image/*">
                
                <input type="submit" value="Register" class="hoverButton w-full p-3 text-white text-lg cursor-pointer gradient-animation rounded mb-10">
                <router-link to="/login" class="hoverButton flex w-full mt-8 justify-center items-center p-3 text-black text-lg cursor-pointer gradient-animation rounded">Login instead</router-link>
            </form>
        </div>
        <div v-if="error !== undefined" class="bg-red-200 border border-red-500 text-red-700 p-3 rounded m-4 text-center w-1/4 mx-auto">
            {{ error }}
        </div>
    </div>
    
    
    </body>
  </template>
  