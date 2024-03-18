import { useState } from 'react';
import useGet from '../fetch/useGet';
import Card from '../UI/Card';
import GreyButton from '../UI/GreyButton';
import SmallButton from '../UI/SmallButton';
import styles from './NotifModal.module.css';
import profile from '../assets/profile.svg';

function NotifModal({open, onClose}) {
    const currUserId = localStorage.getItem("user_id");
   const { error, isLoaded: isLoading, data } = useGet(`/group-member?userid=${currUserId}`)
console.log("members data: ", data)
    if (isLoading) return <div>Loading Noti Modal...</div>
    if (error) return <div>Error: {error.message}</div>
    if (!open) return null;

    function handleClick(e) {
       
    }
    function handleConfirm() {

    }
    return (
        <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modalContainer} >

            <div className={styles.close} onClick={onClose} >X</div>

            {/* {data && data.data.map((user) => ( */}
                <div className={styles.container}>
                    <div className={styles.left}>
                    <img className={styles.img} src={profile} alt=''/>
                    </div>
                    <div className={styles.mid}>
                        <div id='4' className={styles.user} onClick={handleClick}>Username sent you a follow request</div>
                    <div className={styles.btn}>
                    <button className={styles.confirmbtn} onClick={handleConfirm}>Join</button>
                        <button className={styles.removebtn}>Remove</button>
                    </div>
                    </div>
                    <div className={styles.right}></div>
                
                </div>

                <div className={styles.container}>
                    <div className={styles.left}>
                    <img className={styles.img} src={profile} alt=''/>
                    </div>
                    <div className={styles.mid}>
                        <div id='4' className={styles.user} onClick={handleClick}>Username wants you to join GroupTitle</div>
                    <div className={styles.btn}>
                        <button className={styles.confirmbtn} onClick={handleConfirm}>Join</button>
                        <button className={styles.removebtn}>Remove</button>
                    </div>
                    </div>
                    <div className={styles.right}></div>
                </div>

                <div className={styles.container}>
                    <div className={styles.left}>
                    <img className={styles.img} src={profile} alt=''/>
                    </div>
                    <div className={styles.mid}>
                        <div id='4' className={styles.user} onClick={handleClick}>Username wants to join GroupTitle</div>
                    <div className={styles.btn}>
                        <button className={styles.confirmbtn} onClick={handleConfirm}>Accept</button>
                        <button className={styles.removebtn}>Remove</button>
                    </div>
                    </div>
                    <div className={styles.right}></div>
                </div>

                <div className={styles.container}>
                    <div className={styles.left}>
                    <img className={styles.img} src={profile} alt=''/>
                    </div>
                    <div className={styles.mid}>
                        <div id='4' className={styles.user} onClick={handleClick}>GroupTitle added a new event</div>
                        <div className={styles.btn}>
                            <button className={styles.confirmbtn} onClick={handleConfirm}>Accept</button>
                            <button className={styles.removebtn}>Remove</button>
                        </div>
                    </div>
                    <div className={styles.right}>
                        <div className={styles.notif}></div>
                    </div>
                </div>

            </div>
        </div>
            

    )
}

export default NotifModal;