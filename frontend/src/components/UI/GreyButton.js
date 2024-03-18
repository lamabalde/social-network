import styles from './GreyButton.module.css'

const GreyButton = (props) => {
    // const classes = `${styles["lg-btn"]} ${styles[props.className] || ""}`;
    const classes = `${styles["btn"]}` + " " + props.className;
    return (
        <>
            <button className={classes} type={props.type} onClick={props.onClick}>{props.children}</button>
        </>
    )
};

export default GreyButton;