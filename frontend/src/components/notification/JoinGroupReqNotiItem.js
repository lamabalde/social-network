import { useContext, useState } from "react";
import SmallButton from "../UI/SmallButton";
import { WebSocketContext } from "../store/websocket-context";
import Avatar from "../UI/Avatar";
import { GroupsContext } from "../store/groups-context";
import styles from './NotificationItem.module.css'
import SmallAvatar from "../UI/SmallAvatar";


const JoinGroupReqNotiItem = (props) => {
    const wsCtx = useContext(WebSocketContext);
    const grpCtx = useContext(GroupsContext);
    const [isVisible, setIsVisible] = useState(true);


    const grp = grpCtx.groups.find((grp) => grp.id === props.groupId);
    console.log("join grp (noti): ", grp);
    const grpTitle = grp["title"];
    console.log("grp title (noti): ", grpTitle);

    const acceptJoinReqHandler = () => {
        setIsVisible(false);
        console.log("request accepted: ");
        const notiReplyPayloadObj = {};
        notiReplyPayloadObj["label"] = "noti";
        notiReplyPayloadObj["id"] = Date.now();
        notiReplyPayloadObj["type"] = "join-req-reply";
        notiReplyPayloadObj["sourceid"] = props.targetId;
        notiReplyPayloadObj["targetid"] = props.srcUser.id;
        notiReplyPayloadObj["groupid"] = grp.id;
        notiReplyPayloadObj["accepted"] = true;
        console.log("gonna send reply (accept) to join req : ", notiReplyPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(notiReplyPayloadObj));
        let notifarr  =JSON.parse(localStorage.getItem("new_notif"))
        for (let i= 0 ; i < notifarr.length; i++) {
            if (notifarr[i].sourceid == props.srcUser.id && notifarr[i].groupid == props.groupId && notifarr[i].type == "join-req"){
                notifarr.splice(i, 1);
                localStorage.setItem("new_notif", JSON.stringify(Object.values(notifarr)) )
                break
            }
        }
    };
    const declineJoinReqHandler = () => {
        setIsVisible(false);
        console.log("request declined: ");
        const notiReplyPayloadObj = {};
        notiReplyPayloadObj["label"] = "noti";
        notiReplyPayloadObj["id"] = Date.now();
        notiReplyPayloadObj["type"] = "join-req-reply";
        notiReplyPayloadObj["sourceid"] = props.targetId;
        notiReplyPayloadObj["targetid"] = props.srcUser.id;
        notiReplyPayloadObj["groupid"] = grp.id;
        notiReplyPayloadObj["accepted"] = false;
        console.log("gonna send reply (decline) to join req : ", notiReplyPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(notiReplyPayloadObj));
        let notifarr  =JSON.parse(localStorage.getItem("new_notif"))
        for (let i= 0 ; i < notifarr.length; i++) {
            if (notifarr[i].sourceid == props.srcUser.id && notifarr[i].groupid == props.groupId && notifarr[i].type == "join-req"){
                notifarr.splice(i, 1);
                localStorage.setItem("new_notif", JSON.stringify(Object.values(notifarr)) )
                break
            }
        }
    };

    console.log("props.groupid (join)", props.grouptitle);

    return (
        <div>
            {isVisible && (
                <div className={styles.container}>
                    <div className={styles.left}>

                        <SmallAvatar height={50} width={50}/>
                    </div>
                    <div className={styles.mid}>

                        <div>{`${props.srcUser.fname} ${props.srcUser.lname} wants to join ${grpTitle}`}</div>
                        <div className={styles.btn}>

                            <SmallButton onClick={acceptJoinReqHandler}>Accept</SmallButton>
                            <SmallButton onClick={declineJoinReqHandler}>Decline</SmallButton>
                        </div>
                    </div>
                    <div className={styles.right}>
                        {/* <div className={styles.notif}></div> */}
                    </div>
                </div>
            )}
        </div>
    );
};

export default JoinGroupReqNotiItem;