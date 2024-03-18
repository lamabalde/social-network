import { useContext, useEffect, useState } from "react";
import SmallAvatar from "../../UI/SmallAvatar";
import styles from "./OldMsgItem.module.css";
import { UsersContext } from "../../store/users-context";

const NewMsgItem = (props) => {
    const selfId = +localStorage.getItem("user_id");
    const [self, setSelf]= useState();
    const usersCtx = useContext(UsersContext);

    useEffect(()=> {
        setSelf(props.sourceid === selfId)
    },[props])
    console.log("chatprops",props)

    console.log("new chat item", usersCtx.users);
    console.log("new chat item", props.sourceid);

    return (
        <div className={`${self ? styles["self-msg"] : styles["frd-msg"]}`}>

        {!self &&
            <SmallAvatar src={usersCtx.users[props.sourceid-1].avatar}height={30} width={30}></SmallAvatar>
        }
        <div className={styles.wrapper}>
            {/* <div className={`${self ? styles["self-username"] : styles["frd-username"]}`}>{data.data[0].fname} {data.data[0].lname}</div> */}
            <div className={`${self ? styles["self-username"] : styles["frd-username"]}`}>{usersCtx.users[props.sourceid-1].fname} {usersCtx.users[props.sourceid-1].lname}</div>
            <div className={`${self ? styles["chat-bubble-self"] : styles["chat-bubble-frd"]}`}>
                {props.msg}
            </div>
        </div>
        {/* {self &&
            <SmallAvatar height={30} width={30}></SmallAvatar>
        } */}
    </div>
    );
};

export default NewMsgItem;