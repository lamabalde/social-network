import { navigate, useNavigate } from 'react-router-dom';
import Card from '../UI/Card';
import SmallButton from '../UI/SmallButton';
import classes from './UserEvent.module.css'

function UserEvent(props) {
    const navigate = useNavigate();

    function handleClick(e) {
        let id = e.target.id
        // alert(e.target.id)

        if (id == "") {
            id = e.target.parentElement.id
        }
        if (id == "") {
            id = e.target.parentElement.parentElement.id
        }
        navigate("/groupprofile", { state: { id } })

    }
    const moment = require(`moment`);
    const myDate = moment(props.date)
    let newDate;
    if (myDate.isValid()) {
        const mills = myDate.valueOf();
        newDate = new Intl.DateTimeFormat(`en-GB`, {

            day: 'numeric',
            month: 'short',
            year: '2-digit',
            hour: 'numeric',
            minute: 'numeric',


        }).format(mills)
    }


    return <>

    <Card className={classes.card}>
        {/* let [going ,setGoing] = useState */}
        <div id={props.groupid} onClick={handleClick} className={classes.container}>
            <div className={classes.date}>{newDate} </div>
            <div className={classes.title}>{props.title}</div>
            <div>{props.description}</div>
            {/* <div className={classes.btnWrapper}> */}
                {/* {!going && <div id={props.id} className={classes.btn} onClick={handleGoing}>Going</div>} */}
                {/* <SmallButton>Going</SmallButton> */}
                {/* {!notGoing && <div id={props.id} className={classes.btn} onClick={handleNotGoing}>Not Going</div>} */}
                {/* <SmallButton>Not Going</SmallButton> */}
            {/* </div> */}
        </div>
    </Card>
    </>

}

export default UserEvent;
