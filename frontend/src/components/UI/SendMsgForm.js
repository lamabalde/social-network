import React from "react";
import styles from "./SendMsgForm.module.css";

const SendMsgForm = (props) => {
    const classes = `${styles["send-msg-form"]} ${props.className || ""}`
    return (
        <>
            <form className={classes} onSubmit={props.onSubmit} style={props.style}>{props.children}</form>
        </>
    )
};

export default SendMsgForm;