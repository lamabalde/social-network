import styles from "./NotificationBtn.module.css";

const NotificationBtn = (props) => {
    const classes = `${styles["notification-btn"]}` + " " + props.className;
    return (
        <>
            <button className={classes} type={props.type} onClick={props.onClick}>{props.children}</button>
        </>
    )
};

export default NotificationBtn;