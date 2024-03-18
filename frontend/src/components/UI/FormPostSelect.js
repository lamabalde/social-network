import React from "react";
import styles from "./FormPostSelect.module.css";

const FormPostSelect = (props) => {
    const classes = `${styles["select"]} ${props.className || ""}`
    return (
        <>
            <select className={classes} ref={props.reference}>
                {props.options.map((obj, i) => {
                    return <option key={i} value={obj["value"]}>{obj["text"]}</option>
                })}
            </select>
        </>
    )
};

export default FormPostSelect;