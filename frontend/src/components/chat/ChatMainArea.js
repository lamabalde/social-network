import { useEffect, useState, useContext } from "react";
import AllUserChatItems from "./UserChat/AllUserChatItems";
import AllGroupChatItems from "./GroupChat/AllGroupChatItems";
import styles from "./ChatMainArea.module.css";
// import { AuthContext } from "../store/auth-context";
import Chatbox from "./Chatbox/Chatbox.js";
import UserChatDetailTopBar from "./UserChat/UserChatDetailTopBar";
import GroupChatDetailTopBar from "./GroupChat/GroupChatDetailTopBar";

const ChatMainArea = ({grpChat}) => {
    // console.log("user chat followers in chatarea", followersList);

    const [privChatboxOpen, setPrivChatboxOpen] = useState(false);
    const [grpChatboxOpen, setGrpChatboxOpen] = useState(false);
    const [followingOrPublicUserId, setfollowingOrPublicUserId] = useState(0);
    const [grpId, setGrpId] = useState(0);
    const [chatboxReceivesMsg, setChatboxReceivesMsg] = useState(0);

    // const ctx = useContext(AuthContext);

    const openUserChatboxHandler = (followingOrPublicUserId) => {
        console.log("chatbox open for ", followingOrPublicUserId);
        setPrivChatboxOpen(true);
        setfollowingOrPublicUserId(followingOrPublicUserId);
    };

    const closeUserChatboxHandler = () => {
        console.log("chatbox closed for ", followingOrPublicUserId);
        setPrivChatboxOpen(false);
    };

    const openGrpChatboxHandler = (groupId) => {
        console.log("chatbox open for group: ", groupId);
        setGrpChatboxOpen(true);
        setGrpId(groupId);
    };

    const closeGrpChatboxHandler = () => {
        console.log("chatbox closed for ", grpId);
        setGrpChatboxOpen(false);
    };

    // const receiveNewMsgHandler = (chatboxId) => {
    //     console.log("receiving msg from/in chatbox : ", chatboxId);
    //     setChatboxReceivesMsg(chatboxId);
    // };

    // console.log("loggedIn at UserChatMainArea", ctx.isLoggedIn);
    
    return (
        <div 
        className={styles["list"]}
        style={{height: window.innerHeight -110}}
        >

            {!grpChat && !privChatboxOpen &&
                <AllUserChatItems 
                    onOpenChatbox={openUserChatboxHandler}
                    open={privChatboxOpen}
                    grp={grpChat}
                />}
            {!grpChat && privChatboxOpen &&
                <Chatbox 
                    chatboxId={followingOrPublicUserId} 
                    onCloseChatbox={closeUserChatboxHandler}
                    // onReceiveNewMsg={receiveNewMsgHandler}
                    open={privChatboxOpen}
                    grp={grpChat}
                />
            }
            {grpChat && !grpChatboxOpen && 
                <AllGroupChatItems
                    onOpenChatbox={openGrpChatboxHandler}
                    open={grpChatboxOpen}
                />
            }
            {grpChat && grpChatboxOpen && 
                <Chatbox 
                    chatboxId={grpId} 
                    onCloseChatbox={closeGrpChatboxHandler}
                    // onReceiveNewMsg={receiveNewMsgHandler}
                    grp={grpChat}
                /> 
             
                }
        </div>
    );
};

export default ChatMainArea;