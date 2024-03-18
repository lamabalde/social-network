import styles from "./CreateMsgTextarea.module.css";

const CreateMsgTextarea = (props) => {
    const classes = `${styles["textarea"]} ${props.className || ""}`
    return (
        <div>
            <textarea required
                className={classes} 
                name={props.name} 
                id={props.id} 
                placeholder={props.placeholder} 
                rows={props.rows} 
                onChange={props.onChange} 
                ref={props.reference}
                value={props.value}
                style={props.style}
                >{props.children}</textarea>
        </div>
    )
};

export default CreateMsgTextarea;