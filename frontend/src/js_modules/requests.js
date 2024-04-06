import { LOGIN_INVALID_CREDENTIALS } from "./consts"
import { formToObj } from "./helpers/helpers"
import { postsComponent } from "./reactive_elements/postComponent"

export async function sendLogin(event) {
    const loginObj = await formToObj(event)
    const init = {
        method: "POST",
        body: JSON.stringify(loginObj) 
    }
    try {
        const resp = await fetch("http://localhost:8000/loginUser", init)
        if (resp.ok) {
            const parsedResp = await resp.json()
            document.cookie = "userLoggedIn" + "=" + (parsedResp.SessionID || "")  +"; expires="+ Date.now()+3600 + ";SameSite=None; Secure; path=/";
            window.location.href = "/";
            // add the sessionID from response into cookies, then redirect the url to /home
            return
        } else {
            throw new Error(resp.statusText)
        }
    } catch (err) {
        console.log(err)
        return LOGIN_INVALID_CREDENTIALS
    }
}

export async function sendRegister(event) {
    const form = await formToObj(event)

    const init = {
      method: "POST",
      body: JSON.stringify(form),
      headers: {},
    };
    try {
        const resp = await fetch("http://localhost:8000/registerUser", init)
        if (resp.ok) {
            sendLogin(event) // if registering user went fine, then we log the user in
        } else {
            
            const parsedResp = await resp.json()
            //throw new Error(parsedResp);
            console.log(parsedResp)
            return parsedResp.Errors
            
        }
    } catch (err) {
        console.log(err)
        alert(err)
    }
}

export async function submitPostImg(image, postID) {
    const payload = {image:image,postID:postID}
    const init = {
        
        method:"POST",
        body: JSON.stringify(payload)
    }
    try {
      const resp = await fetch("http://localhost:8000/submitPostImg", init);
      if (resp.ok) {
        const parsed = await resp.json()
        return parsed
      } else {
        const parsedResp = await resp.json();
        throw new Error(parsedResp);
      }
    } catch (err) {
      console.log(err);
      alert(err);
    }
}

export async function submitCommentImg(image, commentID) {
  const payload = { image: image, commentID: commentID };
  const init = {
    method: "POST",
    body: JSON.stringify(payload),
  };
  try {
    const resp = await fetch("http://localhost:8000/submitCommentImg", init);
    if (resp.ok) {
      const parsed = await resp.json();
      return parsed;
    } else {
      const parsedResp = await resp.json();
      throw new Error(parsedResp);
    }
  } catch (err) {
    console.log(err);
    alert(err);
  }
}