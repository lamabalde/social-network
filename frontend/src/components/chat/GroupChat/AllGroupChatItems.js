import GroupChatItem from "./GroupChatItem";
import { useCallback, useContext, useEffect, useState } from "react";
import { JoinedGroupContext } from "../../store/joined-group-context";
import { WebSocketContext } from "../../store/websocket-context";
import styles from "../UserChat/AllUserChatItems.module.css";

const AllGroupChatItems = (props) => {
    const joinedGrpCtx = useContext(JoinedGroupContext);
    const wsCtx = useContext(WebSocketContext);
    console.log("ws in AllGrpChatItems: ",wsCtx.websocket);
    console.log("cur user has joined these groups (AllGrpChatItems): ", joinedGrpCtx.joinedGrps);

    useEffect(() => {
        if (wsCtx.websocket !== null && wsCtx.newGroupMsgsObj) {
            console.log("sourceid  (chatitems)", wsCtx.newGroupMsgsObj.sourceid);
            console.log("grpid  (chatitems)", wsCtx.newGroupMsgsObj.groupid);
            // console.log(followingCtx.followingChat.find((follower) => follower.id === wsCtx.newGroupMsgsObj.sourceid));

            if (joinedGrpCtx.joinedGrps && joinedGrpCtx.joinedGrps.find((grp) => grp.id === wsCtx.newGroupMsgsObj.groupid)) {
                // if Cur user is the sender
                console.log("new Received msg data when chatbox is closed (grp)", wsCtx.newGroupMsgsObj);
                console.log("ws receives msg from when chatbox is closed (grp): ", wsCtx.newGroupMsgsObj.groupid);
                wsCtx.newGroupMsgsObj !== null && wsCtx.setNewGroupMsgsObj(null);
                joinedGrpCtx.receiveMsgGroup(wsCtx.newGroupMsgsObj.groupid, false);
            } else {
                console.log("Cur user is not in the group");
            }
        }
    }, [joinedGrpCtx.joinedGrps, wsCtx.newGroupMsgsObj]);

    // middleman, passing grpId from child to parent
    const openGroupChatboxHandler = useCallback((grpId) => props.onOpenChatbox(grpId), [props.onOpenChatbox]);

    return (
        <>
        <div ><h3 className={styles["description"]}>Groups You Have Joined:</h3></div>
        <div>
            {joinedGrpCtx.joinedGrps && joinedGrpCtx.joinedGrps.map((joinedGrp) => {
                // {if (curUserId !== joinedGrp.id) {
                    return <GroupChatItem 
                    key={joinedGrp.id}
                    id={joinedGrp.id}
                    title={joinedGrp.title}
                    creator={joinedGrp.creator}
                    description={joinedGrp.description}
                    // img={group.img}              
                    noti={joinedGrp.chat_noti}
                    onOpenChatbox={openGroupChatboxHandler}
                />
            // }}
            })}
        </div>
        </>
    );
};

export default AllGroupChatItems;