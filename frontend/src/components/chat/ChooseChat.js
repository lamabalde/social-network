import React, { useState } from "react";
import SmallButton from "../UI/SmallButton";
import ChatMainArea from "./ChatMainArea";
import styles from "./ChooseChat.module.css";

const ChooseChat = (props) => {
    const [grpActive, setGrpActive] = useState(false);

    const showUserListHandler = () => {
        console.log("User list");
        setGrpActive(false);
    };

    const showGrpListHandler = () => {
        console.log("Grp list");
        setGrpActive(true);
    };

    return (
        <>
        <div className={styles["switch-bar"]}>
            <SmallButton 
                className={`${!grpActive && styles["active"]} ${styles["switch-bar-btn"]}`}
                onClick={showUserListHandler}
            >Users</SmallButton>

            <SmallButton 
                className={`${grpActive && styles["active"]} ${styles["switch-bar-btn"]}`}
                onClick={showGrpListHandler}
            >Groups</SmallButton>
        </div>
            <ChatMainArea grpChat={grpActive}/>
        </>
    );
};

export default React.memo(ChooseChat);