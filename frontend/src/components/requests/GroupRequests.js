import profile from '../assets/profile.svg';
import SmallButton from '../UI/SmallButton';
import classes from './FollowRequest.module.css';


function GroupRequest() {

    return<div className={classes.container}>
        <div className={classes.wrapper}>
            <img className={classes.img} src={profile}/>
            <div className={classes.user}>Username wants to join Groupname</div>
        </div>
    
        <div>
            <SmallButton className={classes.btn}>Confirm</SmallButton>
        </div>
    </div>
}

export default GroupRequest;