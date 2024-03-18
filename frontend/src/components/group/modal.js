import { useState, useContext, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import useGet from '../fetch/useGet';
import Card from '../UI/Card';
import SmallButton from '../UI/SmallButton';
import styles from './modal.module.css';
import { JoinedGroupContext } from "../store/joined-group-context";

function Modal({ open, onClose, onInvite, currentlyJoined, invitedToJoin }) {
    const currUserId = localStorage.getItem("user_id");
    const jGrpCtx = useContext(JoinedGroupContext);
    const { state } = useLocation();
    const { id } = state; // group id
    const { error, isLoaded: isLoading, data } = useGet(`/group-member?groupid=${id}`)
    if (isLoading) return <div>Loading Modal...</div>
    if (error) return <div>Error: {error.message}</div>
    if (!open) return null;
    
    function inviteHandler(e) {
        window.location.reload();
        const uid = e.target.id;
        jGrpCtx.InviteToJoin(+id, +uid);
        // onInvite(true);
        // setInvitedToJoin(true);
        // setCurrentlyJoined(false);
    }
    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modalContainer} >
                <div className={styles.close} onClick={onClose} >X</div>
                <div className={styles.container}>

                    {data && data.notmembers && data.notmembers.map((user) => (

                        <div key={user.id} id={user.id} className={styles.userContainer}>
                            <div className={styles.img}></div>
                            <div>{user.fname}{user.lname}</div>

                            <div className={styles.end}>
                                {!invitedToJoin && <div className={styles.btn} id={user.id} onClick={inviteHandler}>Send Invitation</div>}
                                {invitedToJoin && <div className={styles.btn} id={user.id} >Invited</div>}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

        </div>

    )
}

export default Modal;



