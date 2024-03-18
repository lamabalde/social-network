import styles from "./FormLabel.module.css";

const FormLabel = (props) => {
    const classes = `${styles["label"]} ${props.className || ""}`
    return (
        <div className={styles["label-container"]}>
            <label className={classes} htmlFor={props.htmlFor}>{props.children}</label>
        </div>
    )
};

export default FormLabel;