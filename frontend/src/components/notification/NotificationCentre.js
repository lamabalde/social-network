import { useState } from "react";
import AllNotificationItems from "./AllNotificationItems";
import styles from "./NotificationCentre.module.css";


const NotificationCentre = (props) => {

    const [notiArr, setNotiArr] = useState([]);
    const selfId = +localStorage.getItem("user_id");

    console.log("noti arr (Notification): ", notiArr);

    return (
        // <div className={styles.overlay} onClick={props.onClose}>
           <>
         {notiArr && 
                 <div className={styles.modalContainer} >
                     <div className={styles.label}>
                         <div>Notifications</div>
                         <div onClick={props.onClose} >X</div>
                     </div>
                     <AllNotificationItems
                        //  notiItems={props.newNoti}
                         onClose={props.onClose}
            
                     />
                     </div>
                 }
                     </>
    );
};

export default NotificationCentre;