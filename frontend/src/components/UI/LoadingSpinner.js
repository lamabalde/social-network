import styles from "./LoadingSpinner.module.css";

const LoadingSpinner = (props) => {
    // const classes = `${styles["loading-spinner"]} ${styles[props.className] || ""}`;
    const classes = `${styles["loading-spinner"]} ${props.className ? styles[props.className] : ""}`;
    return (
        <div className={classes}>
            <div className={styles["spinner"]}>
                <div className={`${styles["dot"]} ${styles["dot-1"]}`}></div>
                {/* <div className={`${styles["dot"]} ${styles["dot-2"]}`}></div> */}
                <div className={`${styles["dot"]} ${styles["dot-3"]}`}></div>
                <div className={`${styles["dot"]} ${styles["dot-4"]}`}></div>
                {/* <div className={`${styles["dot"]} ${styles["dot-5"]}`}></div> */}
                <div className={`${styles["dot"]} ${styles["dot-6"]}`}></div>
                <div className={`${styles["dot"]} ${styles["dot-7"]}`}></div>
                {/* <div className={`${styles["dot"]} ${styles["dot-8"]}`}></div> */}
                <div className={`${styles["dot"]} ${styles["dot-9"]}`}></div>
                <div className={`${styles["dot"]} ${styles["dot-10"]}`}></div>
                {/* <div className={`${styles["dot"]} ${styles["dot-11"]}`}></div> */}
                <div className={`${styles["dot"]} ${styles["dot-12"]}`}></div>
            </div> 
        </div>
    )
};

export default LoadingSpinner;