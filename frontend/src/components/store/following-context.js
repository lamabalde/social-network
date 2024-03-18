import React, { useContext, useEffect, useState } from "react";
import { WebSocketContext } from "./websocket-context";
import { UsersContext } from "./users-context";
export const FollowingContext = React.createContext({
    following: [],
    setFollowing: () => {},
    followingChat: [],
    setFollowingChat: () => {},
    getFollowing: () => {},
    requestToFollow: (followUser) => {},
    follow: (followUser) => {},
    unfollow: (unfollowUser) => {},
    receiveMsgFollowing: (frdId, open, isFollowing) => {},
    // publicChatUsers: [],
    // setPublicChatUsers: () => {},
    otherListedChatUsers: [],
    setOtherListedChatUsers: () => {},
    // chatNotiUserArr: [],
    // setChatNotiUserArr: () => {},
    chatboxOpenedFollowing: (frdId, isFollowing) => {},
});

export const FollowingContextProvider = (props) => {
    const selfId = localStorage.getItem("user_id");
    const followingUrl = `http://localhost:8080/user-following?id=${selfId}`;

    const [following, setFollowing] = useState([]);
    const [followingChat, setFollowingChat] = useState([]);
    // const [publicChatUsers, setPublicChatUsers] = useState([]);
    const [otherListedChatUsers, setOtherListedChatUsers] = useState([]);
    // const [chatNotiUserArr, setChatNotiUserArr] = useState([]);
    const wsCtx = useContext(WebSocketContext);
    const usersCtx = useContext(UsersContext);
    // get following from db
    const getFollowingHandler = () => {
        fetch(followingUrl)
        .then(resp => resp.json())
        .then(data => {
            console.log("followingArr (context): ", data);
            let [followingArr] = Object.values(data);
            setFollowing([...new Set(followingArr)]);
            localStorage.setItem("following", JSON.stringify(followingArr));
        })
        .catch(
            err => console.log(err)
        );
    };

    const getPrivateChatHandler = () => {
        // private chat notification list after login
        // no need to sort coz the db is already returning items ordered by last-msg-time 
        fetch(`http://localhost:8080/private-chat-item?id=${selfId}`)
        .then(resp => resp.json())
        .then(data => {
                console.log("pri chat item data", data);

                const [allChatItemArr] = Object.values(data);
                console.log("followuing", following);
                console.log("followuingChat", followingChat);
                console.log("allChatItemArr", allChatItemArr);

                if (!allChatItemArr || allChatItemArr.length === 0) {
                    setFollowingChat([...new Set(following)]);
                    // filter out following
                    // console.log("public and not following",usersCtx.users.filter(user => user.public === 1));
                    usersCtx.users && usersCtx.users.length !== 0 && !following && setOtherListedChatUsers(usersCtx.users.filter(user => user.public === 1)); // takes care if cur user is not following any user
                    usersCtx.users && usersCtx.users.length !== 0 && following && setOtherListedChatUsers(usersCtx.users.filter(user => user.public === 1 && !following.some(followingUser => followingUser.id === user.id)));

                    console.log("no data", data);
                    return;
                }

                // filter following
                if (allChatItemArr && allChatItemArr.length !== 0) {                    
                    if (!following || following.length === 0) {
                        console.log("no following, have chat item");
                        setFollowingChat([]);
                        // set OtherListedChatUsers to those with chatItems, and public users below them
                        const allOtherChatItems = allChatItemArr.map(chatItem => {
                            const matchedOtherChatItem = usersCtx.users.find(user => user.id === chatItem.sourceid);
                            return { ...chatItem, ...matchedOtherChatItem };
                        });
                        const allPublicUsers = usersCtx.users.filter(user => user.public === 1);
                        // filter out all public users with chat record
                        const filteredPublicUsers = allPublicUsers.filter(user => !allOtherChatItems.some(chatItem => user.id === chatItem.sourceid));
                        setOtherListedChatUsers([...allOtherChatItems, ...filteredPublicUsers]);                        
                    } else {
                        // filter following WITH chat item
                        const filteredFollowingChatItems = allChatItemArr.filter(chatItem => {
                            // console.log("chatItem.sourceid", chatItem.sourceid);
                            if (!following) return false;
                            return following.some(followingUser => followingUser.id === chatItem.sourceid)
                        });
                        console.log("filteredFollowingChatItems", filteredFollowingChatItems);

                        // merge the properties for following WITH chat item
                        const followingChatItems = filteredFollowingChatItems.map(chatItem => {
                            console.log("chatItem.sourceid", chatItem.sourceid);
                            const matchedFollowing = following.find(followingUser => followingUser.id === chatItem.sourceid);
                            return {...chatItem, ...matchedFollowing};
                        });
                        console.log("followingChatItems", followingChatItems);
    
                        // Also display following that are WITHOUT chat item below the ones that have chat item // buggy
                        const filteredFollowingNoChatItems = following.filter(followingUser => {
                            return !allChatItemArr.some(chatItem => followingUser.id === chatItem.sourceid);
                        });
                        console.log("filteredFollowing Without oChatItems", filteredFollowingNoChatItems);

                        let finalFollowingChatItems;
                        if (followingChatItems && followingChatItems.length !== 0 && filteredFollowingNoChatItems &&  filteredFollowingNoChatItems.length !== 0) {
                            finalFollowingChatItems = [...followingChatItems, ...filteredFollowingNoChatItems];
                        } else if (followingChatItems && followingChatItems.length !== 0 && (!filteredFollowingNoChatItems ||  filteredFollowingNoChatItems.length === 0)) {
                            finalFollowingChatItems = [...followingChatItems];
                        } else if ((!followingChatItems || followingChatItems.length === 0) && filteredFollowingNoChatItems && filteredFollowingNoChatItems.length !== 0) {
                            finalFollowingChatItems = [...filteredFollowingNoChatItems];
                        } else {
                            finalFollowingChatItems = [];
                        }
                        setFollowingChat([...new Set(finalFollowingChatItems)]);
    
                        // filter out following, to get all OtherListedChatUsers
                        // first get all OtherListedChatUsers WITH chatItems
                        const filteredOtherChatItems = allChatItemArr.filter(chatItem => {
                            console.log("chatItem.sourceid", chatItem.sourceid);
                            if (!following) return true;
                            return following.every(followingUser => followingUser.id !== chatItem.sourceid);
                        });
                        console.log("filteredOtherChatItems", filteredOtherChatItems);
                        // merge the properties
                        const allOtherChatItems = filteredOtherChatItems.map(chatItem => {
                            console.log("chatItem.sourceid", chatItem.sourceid);
                            const matchedOtherChatItem = usersCtx.users.find(user => user.id === chatItem.sourceid);
                            return {...chatItem, ...matchedOtherChatItem};
                        });
    
                        // case public users WITHOUT chat item
                        let filteredOtherNoChatItems;
                        let allPublicUsersExcludeFolloiwing;
                        if (following && following.length !== 0) allPublicUsersExcludeFolloiwing = usersCtx.users.filter(user => user.public === 1 && !following.some(followingUser => followingUser.id === user.id));
                        console.log("public users: ", allPublicUsersExcludeFolloiwing)
                        if (allPublicUsersExcludeFolloiwing && allPublicUsersExcludeFolloiwing.length !== 0) {
                            filteredOtherNoChatItems = allPublicUsersExcludeFolloiwing.filter(publicUser => {
                                if (!allChatItemArr) return false;
                                return !allChatItemArr.some(chatItem => publicUser.id === chatItem.sourceid);
                            });
                        }
                        console.log("allOtherChatItems", allOtherChatItems);
                        console.log("filteredOther Without ChatItems", filteredOtherNoChatItems);

                        // combine otherUsers WITH chatItems, and public users WITHOUT chatItems to create the finalOtherUsersChatItems
                        let finalOtherUsersChatItems;
                        if (filteredOtherNoChatItems && filteredOtherNoChatItems.length !== 0 && allOtherChatItems && allOtherChatItems.length !== 0) {
                            finalOtherUsersChatItems = [...allOtherChatItems, ...filteredOtherNoChatItems];
                        } else if (filteredOtherNoChatItems && filteredOtherNoChatItems.length !== 0 && (!allOtherChatItems || allOtherChatItems.length === 0)) {
                            finalOtherUsersChatItems = [...filteredOtherNoChatItems];
                        } else if ((!filteredOtherNoChatItems || filteredOtherNoChatItems.length === 0) && allOtherChatItems && allOtherChatItems.length !== 0) {
                            finalOtherUsersChatItems = [...allOtherChatItems];
                        } else {
                            finalOtherUsersChatItems = [];
                        }                                                     
                        console.log("finalOtherNoChatItems", finalOtherUsersChatItems);                  
                        setOtherListedChatUsers(finalOtherUsersChatItems);
                    }
                } 
        }).catch(err => {
            console.log(err);
        })
    };

    const requestToFollowHandler = (followUser) => {
        console.log("request to follow (context): ", followUser.id);

        const followPayloadObj = {};
        followPayloadObj["label"] = "noti";
        followPayloadObj["id"] = Date.now();
        followPayloadObj["type"] = "follow-req";
        followPayloadObj["sourceid"] = +selfId;
        followPayloadObj["targetid"] = followUser.id;
        followPayloadObj["createdat"] = Date.now().toString();
        console.log("gonna send fol req : ", followPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(followPayloadObj));
    };

    const followHandler = (followUser) => {
        if (following && following.length !== 0) { // not empty
            setFollowing(prevFollowing => [...new Set(prevFollowing), followUser]);
            followUser["chat_noti"] = false; // add noti to followUser
            setFollowingChat(prevFollowingChat => [...new Set(prevFollowingChat), followUser]);

            const storedFollowing = JSON.parse(localStorage.getItem("following"));
            const curFollowing = [...storedFollowing, followUser];
            localStorage.setItem("following", JSON.stringify(curFollowing));
        } else {
            setFollowing([followUser]);
            followUser["chat_noti"] = false; // add noti to followUser
            setFollowingChat([followUser]);
            localStorage.setItem("following", JSON.stringify([followUser]));
        }
        console.log("locally stored following (fol)", JSON.parse(localStorage.getItem("following")));
    };

    const unfollowHandler = (unfollowUser) => {
        console.log("unfollowUser (folctx)", unfollowUser);
        setFollowing(prevFollowing => prevFollowing.filter((followingUser) => followingUser.id !== unfollowUser.id));
        setFollowingChat(prevFollowingChat => prevFollowingChat.filter((followingChatUser) => followingChatUser.id !== unfollowUser.id));
        const storedFollowing = JSON.parse(localStorage.getItem("following"));
        const curFollowing = storedFollowing.filter((followingUser) => followingUser.id !== unfollowUser.id);
        localStorage.setItem("following", JSON.stringify(curFollowing));
        console.log("locally stored following (unfol)", JSON.parse(localStorage.getItem("following")));
    };

    // receiveMsgHandler is not only for following, but also for public user chat
    const receiveMsgHandler = (friendId, open, isFollowing) => {
        if (isFollowing) {
            const targetUser = following.find(followingUser => followingUser.id === +friendId);
            console.log("target user", targetUser);
            // noti if not open
            if (!open) {
                console.log("chatbox closed, open=", open);
                targetUser["chat_noti"] = true; // set noti field to true to indicate unread
            } else {
                targetUser["chat_noti"] = false;
                console.log("chatbox opened, open=", open);

                const privateChatNotiPayloadObj = {};
                privateChatNotiPayloadObj["label"] = "set-seen-p-chat-noti";
                privateChatNotiPayloadObj["sourceid"] = friendId;
                privateChatNotiPayloadObj["targetid"] = +selfId;

                if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(privateChatNotiPayloadObj));
            }
            // move userId chat item to the top
            setFollowingChat(prevFollowingChat => [targetUser, ...prevFollowingChat.filter(followingUser => followingUser.id !== +friendId)]);
            console.log("after add chat noti target user", targetUser);
        } else { // if one or both of the users is public and can chat coz of that
            const targetUser = usersCtx.users.find(user => user.id === +friendId);
            console.log("target user", targetUser);

            if (!open) {
                console.log("chatbox closed, open=", open);
                targetUser["chat_noti"] = true; // set noti field to true to indicate unread
            } else {
                targetUser["chat_noti"] = false;
                console.log("chatbox opened, open=", open);
                // delete private chat notification from database
                const privateChatNotiPayloadObj = {};
                privateChatNotiPayloadObj["label"] = "set-seen-p-chat-noti";
                privateChatNotiPayloadObj["sourceid"] = friendId;
                privateChatNotiPayloadObj["targetid"] = +selfId;

                if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(privateChatNotiPayloadObj));
            }
            setOtherListedChatUsers(prevList => [targetUser, ...prevList.filter(otherChatUser => otherChatUser.id !== +friendId)]);
            console.log("after add chat noti target user", targetUser);
        }
    };

    const openFollowingChatboxHandler = (friendId, isFollowing) => {
        let targetUser = null;
        // separated coz the chatlist is reading chat_noti from following for following
        // and users for public or other users
        if (isFollowing) {
            targetUser = following.find(followingUser => followingUser.id === +friendId);
        } else {
            targetUser = usersCtx.users.find(user => user.id === +friendId);
        }
        
        console.log("target user open box", targetUser);

        targetUser["chat_noti"] = false;

        const privateChatNotiPayloadObj = {};
        privateChatNotiPayloadObj["label"] = "set-seen-p-chat-noti";
        privateChatNotiPayloadObj["sourceid"] = friendId;
        privateChatNotiPayloadObj["targetid"] = +selfId;

        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(privateChatNotiPayloadObj));
    };

    useEffect(() => {
        getFollowingHandler();
        getPrivateChatHandler();
        // if (following) {
            // temp list for testing
            // usersCtx.users && setOtherListedChatUsers(usersCtx.users.filter((user) => user.public === 1));
        // }
    // }, [following]);
    }, [usersCtx.users]);

    return (
        <FollowingContext.Provider value={{
            following: following,
            setFollowing: setFollowing,
            followingChat: followingChat,
            setFollowingChat: setFollowingChat,
            getFollowing: getFollowingHandler,
            requestToFollow: requestToFollowHandler,
            follow: followHandler,
            unfollow: unfollowHandler,
            receiveMsgFollowing: receiveMsgHandler,
            // publicChatUsers: publicChatUsers,
            // setPublicChatUsers: setPublicChatUsers,
            otherListedChatUsers: otherListedChatUsers,
            setOtherListedChatUsers: setOtherListedChatUsers,
            // chatNotiUserArr: chatNotiUserArr,
            // setChatNotiUserArr: setChatNotiUserArr,
            chatboxOpenedFollowing: openFollowingChatboxHandler,
        }}>
            {props.children}
        </FollowingContext.Provider>
    );
};