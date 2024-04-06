import  ws  from "../js_modules/ws.js"

export default async function checkAuth() {
  const cookie = getCookieValue("userLoggedIn")
  if (cookie != undefined) {
    const authenticated = await checkSessionID()
    if (authenticated) {
      if (ws.connection.readyState == 0) {
        ws.launchWebsocket('websocket');
      }
      // if ws has not been launched already, launchWebsocket
      return true
    }
  } 
  return false // TODO: change to false once db is working
}

function getCookieValue(name) {
  const regex = new RegExp(`(^| )${name}=([^;]+)`)
  const match = document.cookie.match(regex)
  if (match) {
    return match[2]
  }
}

async function checkSessionID() { // jeeees
  const init = {
    method: "GET",
    credentials: 'include' // this will send the cookie in the request
  }
  try {
    const resp = await fetch("http://localhost:8000/checkAuth", init)
    if (resp.ok) {
      return true
    } else {
      throw new Error(resp.statusText);
    }
  } catch (err) {
    console.log(err)
    return false 
  }
}