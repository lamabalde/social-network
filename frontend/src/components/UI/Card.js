import styles from "./Card.module.css";

const Card = (props) => {
    const classes = `${styles["card"]} ${props.className || ""}`;
    return (
        <>
            <div 
                className={classes} 
                style={props.style} 
                ref={props.reference}
                >
                {props.children}
            </div>
        </>
    )
};

export default Card;