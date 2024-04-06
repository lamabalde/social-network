# Social Network Application

This is a Facebook-like social network application built using React for the frontend and Golang for the backend.

It is hosting on https://socialnetwork-z01.netlify.app/
Note: The websocket is working in localhost, but works inconsistently when hosted.

## Features

The application includes the following features:

- User registration and login with sessions and cookies for authentication
- User profiles with public and private settings
- Ability to follow/unfollow other users
- Creation of posts and comments with privacy settings
- Creation of groups with invitations and requests
- Real-time private and group messaging using Websocket
- Notifications for following requests, group invitations and requests, and group events using Websocket
- Dockerized deployment with separate backend and frontend images

## Technologies Used

- Frontend: React (JavaScript framework)
- Backend: Golang (Go programming language)
- Database: SQLite
- WebSocket: Gorilla WebSocket package
- Migration: golang-migrate package
- Authentication: Sessions and cookies
- Containerization: Docker

## Chat rules

- The user should be able to send private messages to users that he/she is following
- The user that the message was sent to, will receive the message instantly, if he/she is following the user that sent the message (shown under "Users You Are Following:") or if the user has a public profile (shown under "Other Users:")
- if the user is not following the user that sent the message and if the user doesn't have a public profile, the message is stored the database, and it will be shown when he/she starts following the user that sent the message or if the user switches his/her profile to public 
- If a public user switch his/her profile back to private, he/she can still continue on the existing conversations (even if he/she didn't send any messages)
- Groups should have a common chat room, so if a user is a member of the group he/she should be able to send and receive messages to this group chat.





## Social Network Application
  #Application 



#landing-header {
    display:flex;
    justify-content: space-between;
    background-color:aliceblue;
    padding:25px 30px;
    height:30px;
}

#landing-label {
    font-size: 24px;
    color: rgb(4, 75, 137);
    font-family: 'Vina Sans', cursive;
    font-size:38px;
    position: relative;
    top:-10px;
}

.landing-button {
    background-color: aliceblue;
    border: none;
    border: 1px solid rgb(179, 179, 179);
    padding:0 20px;
    height:42px;
    margin-right:20px;
    position: relative;
    top:-6px;
    font-size:16px;
    transition: 0.3s;
}

.landing-button:hover {
    background-color: rgb(201, 229, 254);
    transition: 0.3s;
    transform: scale(1.1);
}

#landing-bg {
    position: absolute;
    top:0;
    left:0;
    background-image: url("./sn-img.jpeg");
    opacity: 0.4;
    width:100vw;
    height: calc(100vh - 80px);
    z-index: -1;
}
