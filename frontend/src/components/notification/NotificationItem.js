import { useContext } from "react";
import { UsersContext } from "../store/users-context";
import FollowReqNotiItem from "./FollowReqNotiItem";
import Avatar from "../UI/Avatar";
import SmallButton from "../UI/SmallButton";
import { useNavigate } from "react-router-dom";
import styles from './NotificationItem.module.css'
import profile from "../assets/profileSmall.svg";
import JoinGroupReqNotiItem from "./JoinGroupReqNotiItem";
import InviteToJoinGroupNotiItem from "./InviteToJoinGroupNotiItem";
import EventNotif from "./eventNotif";

const NotificationItem = (props) => {
    const navigate = useNavigate();
    const usersCtx = useContext(UsersContext);
    const sourceUser = usersCtx.users.find((user) => user.id === props.sourceId);
  
    console.log("props.grouptitle (item)", props);
    return (
        <div>
 
            {props.type === "follow-req" && <FollowReqNotiItem 
            srcUser={sourceUser}
            targetId={props.targetId}
            />}
            {props.type === "join-req" && <JoinGroupReqNotiItem 
            srcUser={sourceUser}
            targetId={props.targetId}
            groupId={props.groupId}
            />}
            {props.type === "invitation" && <InviteToJoinGroupNotiItem 
            srcUser={sourceUser}
            targetId={props.targetId}
            groupId={props.groupId}
            createdAt={props.createdAt}
            />}

            {props.type && props.type.includes("event-notif")  && <EventNotif
            groupId={props.groupId}
            type = {props.type}
            />}
        
        </div>
    );
};

export default NotificationItem;