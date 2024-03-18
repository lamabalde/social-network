import styles from "./FormTextarea.module.css";

const FormTextarea = (props) => {
    const classes = `${styles["textarea"]} ${props.className || ""}`
    return (
        <div>
            <textarea className={classes} name={props.name} id={props.id} placeholder={props.placeholder} rows={props.rows} value={props.value} onChange={props.onChange}>{props.children}</textarea>
        </div>
    )
};

export default FormTextarea;