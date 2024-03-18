import styles from "./UserChatItem.module.css";
import Avatar from "../../UI/Avatar";
import Card from "../../UI/Card";

const UserChatItem = (props) => {
	const defaultImagePath = "user.png";
	// const onlineStatus = false;
	const openChatboxHandler = () => {
		console.log("user chat item clicked");
		props.onOpenChatbox(props.id);
	};

	return (
		<div className={styles["item"]} onClick={openChatboxHandler}>
			{Boolean(props.noti) && (
				<div className={`${styles["noti"]} ${styles["active"]}`}></div>
			)}
			{Boolean(!props.noti) && <div className={styles["noti"]}></div>}
			<Avatar
				className={styles["chat-avatar"]}
				id={props.id}
				src={props.avatar}
				alt=""
				showStatus={true}
				height={"50px"}
				width={"50px"}
			/>
			<div>
				<p
					className={styles["details"]}
				>{`${props.fname} ${props.lname} ${props.nname}`}</p>
			</div>
		</div>
	);
};
export default UserChatItem;
