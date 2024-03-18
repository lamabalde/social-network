import { useNavigate } from "react-router-dom";
import Card from "../../UI/Card";
import { GroupsContext } from "../../store/groups-context";
import styles from "../UserChat/ChatDetailTopBar.module.css";
import { useContext } from "react";
const GroupChatDetailTopBar = (props) => {

    const navigate = useNavigate();

    const grpCtx = useContext(GroupsContext);

    function handleClick() {
        const id =  props.groupId
        navigate("/groupprofile", { state: { id } })
    }

    return (
        <Card className={styles["container"]}>
            <div className={styles.chatdetails}>
                <div className={styles.lnk} onClick={handleClick}>
                    {grpCtx.groups.find(grp => grp.id === props.groupId).title}
                </div>
            </div>
        </Card>
    );
};

export default GroupChatDetailTopBar;