import classes from './followerModal.module.css'
import profile from '../assets/profile.svg';
import { useNavigate } from 'react-router-dom';
import SmallAvatar from '../UI/SmallAvatar';


function FollowingModal(props) {

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
                    <div>Following</div>
                    <div className={classes.close} onClick={props.onClose} >X</div>
                </div>
                <div className={classes.container}>
                    {props.following && props.following.map((follow) => (

                        <div className={classes.container} id={follow.id} onClick={handleClick}>
                            <div className={classes.left }id={follow.id}>
                                <SmallAvatar  className={classes.img} id={follow.id} height={30} width = {30} src={follow.avatar} />
                                <div key={follow.id} id={follow.id}  className={classes.user}>{follow.fname}</div>
                            </div>

                            {/* <div className={classes.right}>  </div> */}
                        </div>

                    ))}


                </div>
            </div>

        </div>

    )
}

export default FollowingModal;
