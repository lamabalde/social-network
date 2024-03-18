import { useContext, useEffect, useState } from "react";
import Avatar from "../../UI/Avatar";
import SmallAvatar from "../../UI/SmallAvatar";
import styles from "./OldMsgItem.module.css";
import { UsersContext } from "../../store/users-context";


const OldMsgItem = (props) => {
    const selfId = +localStorage.getItem("user_id");
    const [self, setSelf] = useState();
    const usersCtx = useContext(UsersContext);

    useEffect(() => {
        setSelf(props.sourceid === selfId)
    }, [props])

    let targetUsers;

    const msgUser = usersCtx.users.find(user => user.id === props.sourceid);

    return (
        <div className={`${self ? styles["self-msg"] : styles["frd-msg"]}`}>

            {!self &&
                <SmallAvatar src={msgUser.avatar}height={30} width={30}></SmallAvatar>
            }
            <div className={styles.wrapper}>
                <div className={`${self ? styles["self-username"] : styles["frd-username"]}`}>{msgUser.fname} {msgUser.lname}</div>
                <div className={`${self ? styles["chat-bubble-self"] : styles["chat-bubble-frd"]}`}>
                    {props.msg}
                </div>
            </div>
        </div>
    );
};

export default OldMsgItem;