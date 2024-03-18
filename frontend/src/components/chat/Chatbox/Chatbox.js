import { useEffect, useContext, useState, useCallback } from "react";
import UsersContext from "../../store/users-context";
import { WebSocketContext } from "../../store/websocket-context";
import ChatboxMsgArea from "../Chatbox/ChatboxMsgArea";
import SendMsg from "./SendMsg";
import { FollowingContext } from "../../store/following-context";
import { JoinedGroupContext } from "../../store/joined-group-context";
import UserChatDetailTopBar from "../UserChat/UserChatDetailTopBar";
import GroupChatDetailTopBar from "../GroupChat/GroupChatDetailTopBar";
import styles from "./Chatbox.module.css";

const Chatbox = (props) => {

    console.log("chatbox props", props)

    const userMsgUrl = "http://localhost:8080/user-message";
    const groupMsgUrl = "http://localhost:8080/group-message";

    const [oldMsgData, setOldMsgData] = useState([]);
    const [newMsgsData, setNewMsgs] = useState([]);
    const [justUpdated, setJustUpdated] = useState(false); // if justUpdated, move chatitem to the top

    const selfId = +localStorage.getItem("user_id");
    const frdOrGrpId = props.chatboxId;
    
    !props.grp && console.log("friendId: ", frdOrGrpId);
    props.grp && console.log("groupId: ", frdOrGrpId);

    // const usersCtx = useContext(UsersContext);
    // console.log("chatbox: ", usersCtx.users);
    const followingCtx = useContext(FollowingContext);
    const joinedGrpCtx = useContext(JoinedGroupContext);
    const wsCtx = useContext(WebSocketContext);
    // console.log("ws in Chatbox: ",wsCtx.websocket);
    // const [msg, setMsg] = useState("");

    // if (wsCtx.websocket !== null) wsCtx.websocket.onmessage = (e) => {
    //     console.log("msg event: ", e);
    //     const msgObj = JSON.parse(e.data);
    //     console.log("ws receives msgObj: ", msgObj);
    //     console.log("ws receives msg: ", msgObj.message);
    //     const newReceivedMsgObj = {
    //         id: msgObj.id,
    //         targetid: msgObj.targetid,
    //         sourceid: msgObj.sourceid,
    //         message: msgObj.message,
    //         createdat: msgObj.createdat,
    //     };
    
    if (!props.grp) {
        // followingCtx.following.find((followingUser) => followingUser.id === frdOrGrpId)["chat_noti"] = false;
        // console.log("following (chatbox)", followingCtx.following);

        // remove noti when following user open chatbox
        if (followingCtx.following && followingCtx.following.includes((following) => following.id === props.chatboxId)) {
            followingCtx.following.find((followingUser) => followingUser.id === frdOrGrpId)["chat_noti"] = false;
            console.log("following (chatbox)", followingCtx.following);
        }
    } else {
        // remove noti when user open group chatbox (wip)
        // if (joinedGrpCtx.joinedGrps && joinedGrpCtx.joinedGrps.includes((following) => following.id === props.chatboxId)) {
        //     followingCtx.following.find((followingUser) => followingUser.id === frdOrGrpId)["chat_noti"] = false;
        //     console.log("following (chatbox)", followingCtx.following);
        console.log("group (chatbox)", joinedGrpCtx.joinedGrps);
    }

    useEffect(() => {
        // private chat
        if (wsCtx.websocket !== null && wsCtx.newPrivateMsgsObj) {
            // if the new msg should be shown in this chatbox
            if (wsCtx.newPrivateMsgsObj.sourceid === frdOrGrpId) {
                console.log("new Received msg data when chatbox is open", wsCtx.newPrivateMsgsObj);
                console.log("ws receives msg from when chatbox is open: ", wsCtx.newPrivateMsgsObj.sourceid);
                setNewMsgs((prevNewMsgs) => [...new Set([...prevNewMsgs, wsCtx.newPrivateMsgsObj])]);
            
                if (wsCtx.newPrivateMsgsObj !== null) wsCtx.setNewPrivateMsgsObj(null);

                // if chatboxId is a user that the cur user is following (not chatting coz of public user)
                if (followingCtx.following && followingCtx.following.find((following => following.id === props.chatboxId))) {
                    followingCtx.receiveMsgFollowing(frdOrGrpId, true, true);
                } else { // public
                    followingCtx.receiveMsgFollowing(frdOrGrpId, true, false);
                }
            }
        }       
        setJustUpdated(prev => !prev);
    }, [wsCtx.newPrivateMsgsObj]) 

    useEffect(() => {
        // group chat
        if (wsCtx.websocket !== null && wsCtx.newGroupMsgsObj) {
            // if the new msg should be shown in this chatbox
            if (wsCtx.newGroupMsgsObj.groupid === frdOrGrpId) {
                console.log("new Received msg data when chatbox is open", wsCtx.newGroupMsgsObj);
                console.log("ws receives msg for group when chatbox is open: ", wsCtx.newGroupMsgsObj.groupid);
                setNewMsgs((prevNewMsgs) => [...new Set([...prevNewMsgs, wsCtx.newGroupMsgsObj])]);
            
                if (wsCtx.newGroupMsgsObj !== null) wsCtx.setNewGroupMsgsObj(null);

                joinedGrpCtx.receiveMsgGroup(frdOrGrpId, true);   
            }
        }
        setJustUpdated(prev => !prev);
    }, [wsCtx.newGroupMsgsObj]) 

    useEffect(() => {
        // clear noti if the chatbox is initilly closed, but then opened
        if (!props.grp) {
            // here is why opening a chatbox will move it to the top!!!
            if (followingCtx.following && followingCtx.following.find((following => following.id === props.chatboxId))) {
                followingCtx.chatboxOpenedFollowing(frdOrGrpId, true);
            } else {
                followingCtx.chatboxOpenedFollowing(frdOrGrpId, false);
            }
        } else {
            joinedGrpCtx.chatboxOpenedGroup(frdOrGrpId, true);
        }
    // props.chatboxId is changed when the chatbox is opened
    }, [props.chatboxId]);

    // send msg to ws
    const sendMsgHandler = useCallback((msg) => {
        let chatPayloadObj = {};
        if (!props.grp) {
            chatPayloadObj["label"] = "private";
            chatPayloadObj["targetid"] = frdOrGrpId;
        } else {
            chatPayloadObj["label"] = "group" ;
            chatPayloadObj["groupid"] = frdOrGrpId;
        }      
        chatPayloadObj["id"] = Date.now();
        chatPayloadObj["sourceid"] = selfId;
        chatPayloadObj["message"] = msg;

        const createdatObj = new Date();

        const selfNewMsgObject = {};
        if (!props.grp) {
            selfNewMsgObject["targetid"] = frdOrGrpId;  
        } else {
            selfNewMsgObject["groupid"] = frdOrGrpId;
            selfNewMsgObject["targetid"] = selfId;
        }  
        selfNewMsgObject["id"] = Date.now();
        selfNewMsgObject["sourceid"] = selfId;
        selfNewMsgObject["message"] = msg;
        selfNewMsgObject["createdat"] = createdatObj.toString();
        
        console.log("new self msg data", selfNewMsgObject);
        setNewMsgs((prevNewMsgs) => [...prevNewMsgs, selfNewMsgObject]);

        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(chatPayloadObj));

        // move friendId chat item to top
        if (props.grp) {
            joinedGrpCtx.receiveMsgGroup(frdOrGrpId, true);
        } else {
            if (followingCtx.following && followingCtx.following.find((following => following.id === props.chatboxId))) {
                followingCtx.receiveMsgFollowing(frdOrGrpId, true, true);
            } else {
                followingCtx.receiveMsgFollowing(frdOrGrpId, true, false);
            }
        } 

        setJustUpdated(prev => !prev);
    }, [props.grp, frdOrGrpId, followingCtx.following, wsCtx.websocket, joinedGrpCtx.receiveMsgGroup, followingCtx.receiveMsgFollowing, props.chatboxId]);

    // const scrolledBottom = (scrolled) => {
    //     scrolled && setJustSent(false);
    // };
    // console.log("new msg data (outside)", newMsgsData);

    const closeChatboxHandler = useCallback(() => {
        props.onCloseChatbox();
    }, []);

    // get old msgsdata.data.push()
    // const AllMsgsToAndFrom = [];

    const fetchOldMsg = async (url) => {
        try {
            const resp = await fetch(url);
            const data = await resp.json();
            console.log("old msg data: ", data);
            if (data.data) {
                const [oldMsgArr] = Object.values(data);
                if (props.grp) { // changing the prop name from sourceid to targetid if group
                    for (let i = 0; i < oldMsgArr.length; i++) {
                        if (oldMsgArr[i].hasOwnProperty("sourceid")) {
                            oldMsgArr[i]["targetid"] = oldMsgArr[i]["sourceid"];
                            // delete oldMsgArr[i]["sourceid"]; // need this to id self
                        }
                    }
                    console.log("old msg arr for grp: ", oldMsgArr);
                }
                oldMsgArr.sort((b, a) => Date.parse(b.createdat) - Date.parse(a.createdat));
                console.log("soreted old msg data", oldMsgArr);
                setOldMsgData([...new Set(oldMsgArr)]);
            }
        } catch(err) {
            console.log(`${err} occurred while fetching old msgs`);
        }
    };

    useEffect(() => {
        // fetch old pri msg
        !props.grp && fetchOldMsg(`${userMsgUrl}?targetid=${selfId}&sourceid=${frdOrGrpId}`);
        // fetch old grp msg
        props.grp && fetchOldMsg(`${groupMsgUrl}?id=${frdOrGrpId}`);        
    }, []);

    return (
        <div className={styles["container"]}>
            <button onClick={closeChatboxHandler} className={styles["close-btn"]}>X</button>
            {props.grp && <GroupChatDetailTopBar groupId={props.chatboxId} />} 
             {!props.grp && <UserChatDetailTopBar userId={props.chatboxId} />}
            {/* <ChatboxMsgArea oldMsgItems={oldMsgData} newMsgItems={newMsgsData}/> */}
            <ChatboxMsgArea oldMsgItems={oldMsgData} newMsgItems={newMsgsData} justUpdated={justUpdated}/>
            <SendMsg onSendMsg={sendMsgHandler}/>            
        </div>
    );
};

export default Chatbox;