import styles from "./LogoutButton.module.css";

const LogoutButton = (props) => {
    // const classes = `${styles["lg-btn"]} ${styles[props.className] || ""}`;
    const classes = `${styles["logout-btn"]}` + " " + props.className;
    return (
        <>
            <button className={classes} type={props.type} onClick={props.onClick}>{props.children}</button>
        </>
    )
};

export default LogoutButton;