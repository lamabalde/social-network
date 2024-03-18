import { useState, useContext, useEffect } from "react";
import { WebSocketContext } from "../store/websocket-context";
import styles from "./Avatar.module.css";
import profile from '../assets/profile.svg'

const SmallAvatar = (props) => {
    // const onlineStatus = false; // change this
    const [onlineStatus, setOnlineStatus] = useState(false);

    const wsCtx = useContext(WebSocketContext);

    useEffect(() => {
        console.log("incoming wsCtx.newOnlineStatusObj.onlineuserids", wsCtx.newOnlineStatusObj.onlineuserids);
        if (wsCtx.websocket !== null && wsCtx.newOnlineStatusObj.onlineuserids) {
            if (wsCtx.newOnlineStatusObj.onlineuserids.find((userId) => props.id === userId)) {
                setOnlineStatus(true);
            } else {
                setOnlineStatus(false);
            }
        }
    },[wsCtx.newOnlineStatusObj.onlineuserids]);

    let defaultImagePath = "user.png";
    if (props.avatar!= null){
        defaultImagePath = props.avatar
    }
    const classes = `${styles["avatar"]} ${props.className || ""}`;
    return (
        <div className={styles["small-wrapper"]}>
            {/* {onlineStatus && <div className={styles["small-online-status-dot"]}></div>}
            {!onlineStatus && <div className={styles["small-offline-status-dot"]}></div>} */}
            {props.src && <img className={classes} src={props.src} alt={props.alt} height={props.height} width={props.width}/>}
            {!props.src && <img className={classes} src={require("../../images/"+`${defaultImagePath}`)} alt={props.alt} height={props.height} width={props.width}/>}
        </div>
    )
};

export default SmallAvatar;