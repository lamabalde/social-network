import styles from "../UserChat/UserChatItem.module.css";
import Card from "../../UI/Card";

const GroupChatItem = (props) => {
    const openChatboxHandler = () => {
        console.log("group chat item clicked");
        props.onOpenChatbox(props.id);
    };

    return (
        <div className={styles["item"]} onClick={openChatboxHandler}>
            {Boolean(props.noti) && <div className={`${styles["noti"]} ${styles["active"]}`}></div>}
            {Boolean(!props.noti) && <div className={styles["noti"]}></div>}
            {/* <img className={styles["avatar"]} src={props.img} alt="" height={"50px"} width={"50px"}/> */}
            <div><p className={styles["details"]}>{`${props.title}`}</p></div>
        </div>
    );
};

export default GroupChatItem;