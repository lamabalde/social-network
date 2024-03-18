import React, { useContext, useEffect, useState } from "react";
import { WebSocketContext } from "./websocket-context";
import { GroupsContext } from "./groups-context";

export const JoinedGroupContext = React.createContext({
    joinedGrps: [],
    // requestedGrps: [],
    setJoinedGrps: () => {},
    // joinedGrpsChat: [],
    // setJoinedGrpsChat: () => {},
    getJoinedGrps: () => {},
    requestToJoin: (joinGrpId) => {},
    InviteToJoin: (grp, InvitedUser) => {},
    join: (toJoinGrp, user) => {},
    // leave: (toLeaveGrp, user) => {},
    // storeGroupMember: (userId, groupId) => {},
    receiveMsgGroup: (groupId, open) => {},
    // chatNotiUserArr: [],
    // setChatNotiUserArr: () => {},
    chatboxOpenedGroup:(groupId) => {},
});

export const JoinedGroupContextProvider = (props) => {
    const selfId = +localStorage.getItem("user_id");
    const grpCtx = useContext(GroupsContext);

    const [joinedGrps, setJoinedGrps] = useState([]);
    // const [joinedGrpsChat, setJoinedGrpsChat] = useState([]);
    const [requestedGroups, setRequestedGroups] = useState([]);
    const joinedGroupingUrl = `http://localhost:8080/group-member?userid=${selfId}`;
    const requestedGroupUrl = `http://localhost:8080/group-request-by-user?id=${selfId}`;
    const wsCtx = useContext(WebSocketContext);

    // get from db
    const getJoinedGrpsHandler = () => {
        fetch(joinedGroupingUrl)
        .then(resp => resp.json())
        .then(data => {
            console.log("joinedGroupsArr (context) ", data);
            let [joinedGroupsArr] = Object.values(data);
            setJoinedGrps(joinedGroupsArr);
            // localStorage.setjoinedGrpsItem("joinedGroups", JSON.stringify(joinedGroupsArr));
            localStorage.setItem("joinedGroups", JSON.stringify(joinedGroupsArr));

        })
        .catch(
            err => console.log(err)
        );
    };
    const getRequestedGrpsHandler = () => {
        console.log("RUNNNING")
        fetch(requestedGroupUrl)
        .then(resp => resp.json())
        .then(data => {
            console.log("requestedGroup (context): ", data);
            let [requestedGroupsArr] = Object.values(data);
            setRequestedGroups(requestedGroupsArr);
            localStorage.setItem("requestedGroups", JSON.stringify(requestedGroupsArr));
        })
        .catch(
            err => console.log(err)
        );
    };

    const fetchGroupChatData = async (url) => {
        // no need to sort coz the db is already returning items ordered by last-msg-time
        const resp = await fetch(url);
        const data = await resp.json();
        const [dataArr] = Object.values(data);
        console.log(" data Arr", dataArr);
        return dataArr;
    };

    const getGroupChatHandler = () => {
        (async function() {
            try {
                const promises = [
                    fetchGroupChatData(joinedGroupingUrl),
                    fetchGroupChatData(`http://localhost:8080/group-chat-item?userid=${selfId}`),
                ];
                const result = await Promise.all(promises);

                const joinedGroupsArr = result[0];
                console.log("memberOfGrpArr", joinedGroupsArr);
                if (!joinedGroupsArr) {
                    setJoinedGrps([]);
                    console.log("user hasn't joined any group");
                    return;
                } else {
                    setJoinedGrps(joinedGroupsArr);
                    console.log("user has joined groups", joinedGroupsArr);
                }
                localStorage.setItem("joinedGroups", JSON.stringify(joinedGroupsArr));

                const allGrpChatItemArr = result[1];
                console.log("grp chat item Arr", allGrpChatItemArr);
                if (!allGrpChatItemArr) {
                    console.log("no grp chat record", allGrpChatItemArr);
                } else if (allGrpChatItemArr) {
                    const filteredGroupChatItems = allGrpChatItemArr.filter(chatItem => {
                        // console.log("joinedGroupsArr length", joinedGroupsArr.length);
                        // for (let grp of joinedGroupsArr) {
                        //     console.log("grp id", grp.id);
                        // }
                        console.log("returned", joinedGroupsArr.some(grp => grp.id == chatItem.groupid && +selfId === chatItem.userid));
                        return joinedGroupsArr.some(grp => grp.id === chatItem.groupid && +selfId === chatItem.userid);
                    });
                    console.log("filteredGroupChatItems", filteredGroupChatItems);

                     // merge the properties
                     const groupChatItems = filteredGroupChatItems.map(chatItem => {
                        const matchedGroups = joinedGroupsArr.find(grp => grp.id === chatItem.groupid);
                        return {...chatItem, ...matchedGroups};
                    });
                    //  === chatItem.groupid && +selfId === chatItem.userId
                    // Also display groups even if there is no chat item
                    let filteredGroupsNoChatItems;
                    if (joinedGroupsArr) {
                        filteredGroupsNoChatItems = joinedGroupsArr.filter(grp => {
                            if (!allGrpChatItemArr) return false;
                            return !allGrpChatItemArr.some(chatItem => grp.id);
                        });
                    }
                    console.log("filteredFollowing Without oChatItems", filteredGroupsNoChatItems);
                    const finalGroupChatItems = [...groupChatItems, ...filteredGroupsNoChatItems];
                    setJoinedGrps(finalGroupChatItems);
                }



                 // Also display joined groups even if there is no chat item

            } catch(err) {
                console.log(`${err} occurred during fetch`);
            }
        }());
    };

    const requestToJoinHandler = (joinGrpId) => {
        getRequestedGrpsHandler()
        console.log("request to join user (context): ", +selfId);
        console.log("request to join grp (context): ", joinGrpId);
        const grp = grpCtx.groups.find((grp) => grp.id === joinGrpId);
        console.log("join grp (context): ", grp);
        const creatorId = grp["creator"];
        console.log("creator of join grp (context): ", creatorId);

        const joinGrpPayloadObj = {};
        joinGrpPayloadObj["label"] = "noti";
        joinGrpPayloadObj["id"] = Date.now();
        joinGrpPayloadObj["type"] = "join-req";
        joinGrpPayloadObj["sourceid"] = +selfId;
        joinGrpPayloadObj["targetid"] = creatorId;
        joinGrpPayloadObj["groupid"] = joinGrpId;
        joinGrpPayloadObj["createdat"] = Date.now().toString();
        console.log("gonna send join grp req : ", joinGrpPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(joinGrpPayloadObj));
    };

    const InviteToJoinHandler = (grpid, InvitedUserId) => {
        console.log("Invite to join user (context): ", InvitedUserId);
        console.log("Invite to join grp (context): ", grpid);

        const grp = grpCtx.groups.find((grp) => grp.id === grpid);
        console.log("invite grp (context): ", grp);
        const creatorId = grp["creator"];
        console.log("creator of invite grp (context): ", creatorId);

        const InviteToJoinPayloadObj = {};
        InviteToJoinPayloadObj["label"] = "noti";
        InviteToJoinPayloadObj["id"] = Date.now();
        InviteToJoinPayloadObj["type"] = "invitation";
        InviteToJoinPayloadObj["sourceid"] = creatorId;
        InviteToJoinPayloadObj["targetid"] = InvitedUserId;
        InviteToJoinPayloadObj["groupid"] = grpid;
        InviteToJoinPayloadObj["createdat"] = Date.now().toString();
        console.log("gonna send invite : ", InviteToJoinPayloadObj);
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(InviteToJoinPayloadObj));
    };

    const joinHandler = (toJoinGrp, user) => {
        console.log("toJoinGrp (jg ctx)", toJoinGrp);
        console.log("user (jg ctx)", user);
        if (joinedGrps) { // not empty
            setJoinedGrps(prevJoinedGrps => [...prevJoinedGrps, toJoinGrp]);
            // const storedJoinedGrps = JSON.parse(localStorage.getItem("joined-grps"));
            // const curJoined = [...storedJoinedGrps, toJoinGrp];
            // localStorage.setItem("joined", JSON.stringify(curJoined));
        } else {
            setJoinedGrps([toJoinGrp]);
            localStorage.setItem("joined-grps", JSON.stringify([toJoinGrp]));
        }
        console.log("locally stored joined grp (jg ctx)", JSON.parse(localStorage.getItem("joined-grps")));
    };

    // const leaveHandler = (toLeaveGrp, user) => {
    //     console.log("user (leaveHandler)", user);
    //     console.log("leave grp (leaveHandler)", toLeaveGrp);
    //     setJoinedGrps(prevJoinedGrps => prevJoinedGrps.filter((prevJoinedGrp) => prevJoinedGrp.id !== toLeaveGrp.id));
    //     const storedJoinedGrps = JSON.parse(localStorage.getItem("joined-grps"));
    //     const curJoinedGrps = storedJoinedGrps.filter((prevJoinedGrp) => prevJoinedGrp.id !== toLeaveGrp.id);
    //     localStorage.setItem("joined-grps", JSON.stringify(curJoinedGrps));
    //     console.log("locally stored joined-grps (leaveHandler)", JSON.parse(localStorage.getItem("joined-grps")));
    // };

    // const storeGroupMemberHandler = () => {
        // store the new member to group member db table
    // };

    const receiveMsgHandler = (groupId, open) => {
        const targetGroup = joinedGrps.find(group => group.id === +groupId);
        console.log("target group", targetGroup);
        // noti if not open
        if (!open) {
            console.log("group chatbox closed, open=", open);
            targetGroup["chat_noti"] = true; // set noti field to true to indicate unread
        } else {
            targetGroup["chat_noti"] = false;
            console.log("group chatbox opened, open=", open);

            const groupChatNotiPayloadObj = {};
            groupChatNotiPayloadObj["label"] = "set-seen-g-chat-noti";
            groupChatNotiPayloadObj["sourceid"] = +selfId;
            groupChatNotiPayloadObj["groupid"] = groupId;

            if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(groupChatNotiPayloadObj));
        }
        // move userId chat item to the top
        setJoinedGrps(prevJoinedGrps => [targetGroup, ...prevJoinedGrps.filter(joinedGrp => joinedGrp.id !== +groupId)]);
        console.log("after add chat noti target group", targetGroup);

    };

    const openGroupChatboxHandler = (groupId) => {
        const targetGroup = joinedGrps.find(group => group.id === +groupId);
        console.log("targetGroup open box", targetGroup);

        targetGroup["chat_noti"] = false;

        const groupChatNotiPayloadObj = {};
        groupChatNotiPayloadObj["label"] = "set-seen-g-chat-noti";
        groupChatNotiPayloadObj["sourceid"] = +selfId;
        groupChatNotiPayloadObj["groupid"] = +groupId;

        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(groupChatNotiPayloadObj));
    };

    useEffect(() => {
        // getJoinedGrpsHandler();
        getGroupChatHandler();
    }, []);
    useEffect(() => getRequestedGrpsHandler(), []);

    return (
        <JoinedGroupContext.Provider value={{
            joinedGrps: joinedGrps,
            setJoinedGrps: setJoinedGrps,
            // joinedGrpsChat: joinedGrpsChat,
            // setJoinedGrpsChat, setJoinedGrpsChat,
            getJoinedGrps: getJoinedGrpsHandler, // implement
            requestToJoin: requestToJoinHandler,
            InviteToJoin: InviteToJoinHandler,
            requestLocalStrg: getRequestedGrpsHandler,
            requestedGroups:requestedGroups,
            // getFollowing : getRequestedGrpsHandler,
            join: joinHandler,
            // leave: leaveHandler,
            // storeGroupMember: storeGroupMemberHandler,
            receiveMsgGroup: receiveMsgHandler,
            // chatNotiUserArr: chatNotiUserArr,
            // setChatNotiUserArr: setChatNotiUserArr,
            chatboxOpenedGroup: openGroupChatboxHandler,
        }}>
            {props.children}
        </JoinedGroupContext.Provider>
    );
};