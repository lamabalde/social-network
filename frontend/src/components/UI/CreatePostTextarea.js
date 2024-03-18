import styles from "./CreatePostTextarea.module.css";

const CreatePostTextarea = (props) => {
    const classes = `${styles["textarea"]} ${props.className || ""}`
    return (
        <div>
            <textarea className={classes} name={props.name} id={props.id} placeholder={props.placeholder} value={props.value} rows={props.rows} onChange={props.onChange} ref={props.reference} required>{props.children}</textarea>
        </div>
    )
};

export default CreatePostTextarea;