import styles from './NotificationItem.module.css'
import profile from "../assets/profileSmall.svg";
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import useGet from '../fetch/useGet';
import Avatar from '../UI/Avatar';
import SmallAvatar from '../UI/SmallAvatar';



function EventNotif(props) {
    const navigate = useNavigate();
    const [isVisible, setIsVisible] = useState(true);

    const { error, isLoaded: isLoading, data } = useGet(`/group?id=${props.groupId}`)    
    if (isLoading) return <div>Loading Event Noti...</div>
    if (error) return <div>Error: {error.message}</div>

    function handleClick(e) {
        setIsVisible(false);
        console.log("click")
        const id = e.target.id
        let notifarr  =JSON.parse(localStorage.getItem("new_notif"))
       let newarray =  notifarr.filter((obj) => (obj.groupid != id))
       console.log("newarray ", newarray)


    localStorage.setItem("new_notif", JSON.stringify(Object.values(newarray)) )
        navigate("/groupprofile", { state: { id } })
        console.log("5678", id)
    }

    return (
        <div>
            {isVisible && (

                <div className={styles.container}>
                    <div className={styles.left}>
                        {/* <img className={styles.img} src={profile} alt='' /> */}
                        <SmallAvatar height={50} width={50}/>
                    </div>
                    <div className={styles.mid}>
                        { data && data.data && 
                        <div id={props.groupId} onClick={handleClick} className={styles.user}> {data.data[0].title} added new event: {props.type.split("+")[1]}</div>
                        }
                        {/* <div id={props.groupId} onClick={handleClick} className={styles.user}>GroupTitle {props.groupId} added a new event: EventTitle</div>   */}
                    </div>
                    <div className={styles.right}>
                        {/* <div className={styles.notif}></div> */}
                    </div>
                </div>
            )}
        </div>
    )
}


export default EventNotif;