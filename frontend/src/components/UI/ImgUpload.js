import styles from "./ImgUpload.module.css";
// import img from "../assets/img.svg"

const ImgUpload = (props) => {
    // const classes = `${styles["lg-btn"]} ${styles[props.className] || ""}`;
    const classes = `${styles["input"]}` + " " + props.className;
    return (
        <>
            <label htmlFor={props.id} className={styles["label"]}>{props.text}</label>
            {/* <label  htmlFor={props.id} className={styles["label"]}><img src={img}></img></label> */}

            <input type="file" name={props.name} id={props.id} accept={props.accept} onChange={props.onChange} className={classes}/>
        </>
    )
};

export default ImgUpload;