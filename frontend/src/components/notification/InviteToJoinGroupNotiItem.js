import { useContext, useState } from "react";
import SmallButton from "../UI/SmallButton";
import { WebSocketContext } from "../store/websocket-context";
import Avatar from "../UI/Avatar";
import { GroupsContext } from "../store/groups-context";
import { JoinedGroupContext } from "../store/joined-group-context";
import styles from './NotificationItem.module.css'
import SmallAvatar from "../UI/SmallAvatar";


const InviteToJoinGroupNotiItem = (props) => {
    const [isVisible, setIsVisible] = useState(true);

    const wsCtx = useContext(WebSocketContext);
    const grpCtx = useContext(GroupsContext);
    const jGrpCtx = useContext(JoinedGroupContext)
    const grp = grpCtx.groups.find((grp) => grp.id === props.groupId);
    console.log("join grp (noti): ", grp);
    const grpTitle = grp["title"];
    console.log("grp title (noti): ", grpTitle);

    const acceptInvitationHandler = () => {
        setIsVisible(false);
        console.log("request accepted: ");
        const notiReplyPayloadObj = {};
        notiReplyPayloadObj["label"] = "noti";
        notiReplyPayloadObj["id"] = Date.now();
        notiReplyPayloadObj["type"] = "invitation-reply";
        notiReplyPayloadObj["sourceid"] = props.targetId;
        notiReplyPayloadObj["targetid"] = props.srcUser.id;
        notiReplyPayloadObj["groupid"] = grp.id;
        notiReplyPayloadObj["accepted"] = true;
        console.log("gonna send reply (accept) to Invitation : ", notiReplyPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(notiReplyPayloadObj));
        jGrpCtx.getFollowing();

        let notifarr  =JSON.parse(localStorage.getItem("new_notif"))
        // let newarray =  notifarr.filter((obj) => (obj.sourceid != props.srcUser.id && obj.groupid != props.groupId && obj.type != "invitation"))
        for (let i= 0 ; i < notifarr.length; i++) {
            if (notifarr[i].sourceid == props.srcUser.id && notifarr[i].groupid == props.groupId && notifarr[i].type == "invitation"){
            notifarr.splice(i, 1);
                localStorage.setItem("new_notif", JSON.stringify(Object.values(notifarr)) )
                break
            }
        }
        
    };
    const declineInvitationHandler = () => {
        setIsVisible(false);

        console.log("request declined: ");
        const notiReplyPayloadObj = {};
        notiReplyPayloadObj["label"] = "noti";
        notiReplyPayloadObj["id"] = Date.now();
        notiReplyPayloadObj["type"] = "invitation-reply";
        notiReplyPayloadObj["sourceid"] = props.targetId;
        notiReplyPayloadObj["targetid"] = props.srcUser.id;
        notiReplyPayloadObj["groupid"] = grp.id;
        notiReplyPayloadObj["accepted"] = false;
        console.log("gonna send reply (decline) to Invitation : ", notiReplyPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(notiReplyPayloadObj));
        let notifarr  =JSON.parse(localStorage.getItem("new_notif"))
        for (let i= 0 ; i < notifarr.length; i++) {
            if (notifarr[i].sourceid == props.srcUser.id && notifarr[i].groupid == props.groupId && notifarr[i].type == "invitation"){
                notifarr.splice(i, 1);
                localStorage.setItem("new_notif", JSON.stringify(Object.values(notifarr)) )
                break
            }
        }
    };

    return (
        <div>
            {isVisible && (

                <div className={styles.container}>
                    <div className={styles.left}>

                        <SmallAvatar height={50} width={50}/>
                    </div>
                    <div className={styles.mid}>

                        <div>{`${props.srcUser.fname} ${props.srcUser.lname} wants you to join ${grpTitle}`}</div>
                        <div className={styles.btn}>

                            <SmallButton onClick={acceptInvitationHandler}>Accept</SmallButton>
                            <SmallButton onClick={declineInvitationHandler}>Decline</SmallButton>
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


export default InviteToJoinGroupNotiItem;