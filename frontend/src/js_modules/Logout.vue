<script setup>
import { routerKey, useRouter } from 'vue-router';
import { currentUser } from './reactive_elements/userComponent';
import ws from './ws';

async function sendLogout() {
    const init = {
        method: "POST",
        credentials: 'include', // this will send the cookie in the request
    }
    try {
        const resp = await fetch("http://localhost:8000/logout", init)
        if (resp.ok) {
            
            ws.closeWebsocket()
            document.cookie = "userLoggedIn=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"; 
            window.location.href = "/login";
            currentUser.value.loaded = false
            console.log("ws readyState: ", ws.connection.readyState)
            console.log("logout ok");
            return
        } else {
            
            const parsedResp = await resp.json()
            throw new Error(parsedResp)
        }
    } catch (err) {
        console.log(err)
        alert(err)
    }
}
sendLogout()


</script>

<template>


</template>