import { useState, useContext, useEffect } from "react";
import { WebSocketContext } from "../store/websocket-context";
import styles from "./Avatar.module.css";
import profile from '../assets/profile.svg'

const Avatar = ({className, id, src, alt, height, width, showStatus}) => {
    // const onlineStatus = false; // change this
    const [onlineStatus, setOnlineStatus] = useState(false);

    const wsCtx = useContext(WebSocketContext);

    useEffect(() => {
        if (wsCtx.websocket !== null && wsCtx.newOnlineStatusObj.onlineuserids) {
            if (wsCtx.newOnlineStatusObj.onlineuserids.find((userId) => id === userId)) {
                setOnlineStatus(true);
            } else {
                setOnlineStatus(false);
            }
        }
    },[wsCtx.newOnlineStatusObj.onlineuserids]);

    const defaultImagePath = "user.png";
    const classes = `${styles["avatar"]} ${className || ""}`;
    return (
        <div className={styles["wrapper"]}>
            {showStatus && onlineStatus && <div className={styles["online-status-dot"]}></div>}
            {showStatus && !onlineStatus && <div className={styles["offline-status-dot"]}></div>}
            {src && <img className={classes} src={src} alt={alt} height={height} width={width}/>}
            {!src && <img className={classes} src={require("../../images/"+`${defaultImagePath}`)} alt={alt} height={height} width={width}/>}
        </div>
    )
};

export default Avatar;