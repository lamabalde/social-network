import classes from './followerModal.module.css'
import { useNavigate } from 'react-router-dom';
import SmallAvatar from '../UI/SmallAvatar';


function FollowerModal(props) {
    const navigate = useNavigate();


    function handleClick(e) {
        const id = e.target.id

        console.log("id: ", id)
        navigate(`/profile/${id}`)

    }

    return (
        <div className={classes.overlay} onClick={props.onClose}>
            <div className={classes.modalContainer} >
                <div className={classes.label}>
                    <div></div>
                    <div>Followers</div>
                    <div className={classes.close} onClick={props.onClose} >X</div>
                </div>
                <div className={classes.container}>
                    {props.followers && props.followers.map((follower) => (

                        <div className={classes.container} id={follower.id} onClick={handleClick} >
                            <div className={classes.left} id={follower.id}>
                                <SmallAvatar  className={classes.img} id={follower.id} height={30} width = {30} src={follower.avatar} />

                                <div key={follower.id} id={follower.id} className={classes.user}>{follower.fname}</div>
                            </div>

                            {/* <div className={classes.right}>  </div> */}
                        </div>

                    ))}


                </div>
            </div>

        </div>

    )
}

export default FollowerModal;
