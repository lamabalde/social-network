import classes from './FollowRequest.module.css';
import profile from '../assets/profile.svg';
import SmallButton from '../UI/SmallButton';
import { useNavigate } from 'react-router-dom';

function FollowRequest() {
    const navigate = useNavigate();

    function handleClick(e) {
        let id = e.target.id  

        console.log("follow id",id);
        navigate("/profile", {state: { id } })
    }

    function handleClick() {
        
    }
    return<div className={classes.container}>
        <div className={classes.wrapper}>
            <img className={classes.img} src={profile} alt=''/>
            <div id='4' className={classes.user} onClick={handleClick}>Username</div>
        </div>
    
        <div>
            <SmallButton className={classes.btn} onClick={handleClick}>Confirm</SmallButton>
        </div>
    </div>
}

export default FollowRequest;
