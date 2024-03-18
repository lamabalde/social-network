import styles from "./FormInput.module.css";

const FormInput = (props) => {
    const classes = `${styles["input"]} ${props.className || ""}`;
    
    return (
        <div className={styles["input-container"]}>
           <input required className={classes}
                    onChange={props.onChange}
                    id={props.id}
                    type={props.type}
                    name={props.name}
                    value={props.value}
                    placeholder={props.placeholder}
                    accept={props.accept}
                    />
        </div>
    )
};

export default FormInput;