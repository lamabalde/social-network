import { useCallback, useContext } from "react";
import Avatar from "../../UI/Avatar";
import Card from "../../UI/Card";
import { UsersContext } from "../../store/users-context";
import styles from "./ChatDetailTopBar.module.css";
import { Link, useNavigate } from "react-router-dom";
import SmallAvatar from "../../UI/SmallAvatar";

const UserChatDetailTopBar = (props) => {
    const usersCtx = useContext(UsersContext);
    console.log("tar uid", props.userId);
    const targetUser = usersCtx.users.find(user => +user.id === props.userId);
    console.log("tar u", targetUser);

    return (
        <Card className={styles["container"]}>
            <div className={styles.chatdetails}>
                <SmallAvatar src={targetUser.avatar}height={40} width={40} />
                <Link to={`profile/${props.userId}`} className={styles.lnk}>{targetUser.fname} {targetUser.lname}</Link>
            </div>
        </Card>
    );
};

export default UserChatDetailTopBar;