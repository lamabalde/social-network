import { useEffect, useState } from "react";
import NotificationItem from "./NotificationItem";


const AllNotificationItems = (props) => {

    console.log("all notficifjw", props)

    const storedNotif = JSON.parse(localStorage.getItem("new_notif"));
    const [notiArr, setNotiArr] = useState([]);

    useEffect(() => {
        // const storedNotif = JSON.parse(localStorage.getItem("new_notif"));
        setNotiArr(storedNotif)
    }, []);

    useEffect(() => {

        localStorage.setItem("new_notif", JSON.stringify(Object.values(notiArr)))

    }, [notiArr]);

    console.log("last exit before bridge: ", notiArr)
    return (
        <div>
            {notiArr && notiArr.map((notiItem) => {
                return (
                    <NotificationItem
                        key={notiItem.id}
                        id={notiItem.id}
                        type={notiItem.type}
                        targetId={notiItem.targetid}
                        sourceId={notiItem.sourceid}
                        createdAt={notiItem.createdat}
                        groupId={notiItem.groupid}
                    />
                );
            })
            }
        </div>
    );
};

export default AllNotificationItems;