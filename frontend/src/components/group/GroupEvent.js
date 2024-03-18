import { useEffect, useState } from 'react';
import Card from '../UI/Card';
import GreyButton from '../UI/GreyButton';
import SmallButton from '../UI/SmallButton';
import classes from './GroupEvent.module.css';


function GroupEvent(props) {
    let userid = localStorage.getItem("user_id")

    const [going, setGoing] = useState(false)
    const [notGoing, setNotGoing] = useState(false)
    const [eventMemberData, setEventMemberData] = useState([])

    useEffect(() => {

        fetch(`http://localhost:8080/group-event-member?id=${props.id}`)
            .then(resp => resp.json())
            .then(data => {
                setEventMemberData(data.data)

            })
            .catch(
                err => console.log(err)
            );


    }, []);

    useEffect(() => {
        eventMemberData && eventMemberData.map((member) => {
            console.log({member})
            if (member.userid === (parseInt(userid))) {
                if (member.status == 2) {
                    setGoing(true)
                    setNotGoing(false)

            } else if (member.status == 3) {
                setNotGoing(true)
                setGoing(false)

            } else {
                setNotGoing(false)
                setGoing(false)
            }
        }
        })

    }, [eventMemberData]);


    // var myDate = new Date(props.date);
    // var options = {
    //     day: '2-digit',
    //     month: 'short',
    //     year: '2-digit',
    //     hour: '2-digit',
    //     minute: '2-digit',
    //     hour12: false
    //   };
    //   const newDate = myDate.toLocaleString("en-GB", options)
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

    const currUserId = localStorage.getItem("user_id");

    function handleNotGoing(e) {
        setNotGoing(true)
        setGoing(false)
        const id = e.target.id;
        const data = {
            id: 0,
            status: 3,
            userid: parseInt(currUserId),
            eventid: parseInt(id),
        };

        fetch('http://localhost:8080/group-event-member',
            {
                method: 'POST',
                credentials: "include",
                mode: "cors",
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(() => {
                // navigate.replace('/??')
                console.log("posted")
            })
    }

    function handleGoing(e) {
        setGoing(true)
        setNotGoing(false)
        const id = e.target.id;

        const data = {
            id: 0,
            status: 2,
            userid: parseInt(currUserId),
            eventid: parseInt(id),
        };

        fetch('http://localhost:8080/group-event-member',
            {
                method: 'POST',
                credentials: "include",
                mode: "cors",
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(() => {
                // navigate.replace('/??')
                console.log("posted")
            })
    }

    return <Card className={classes.card}>

        <div className={classes.container}>
            <div className={classes.date}>{newDate} </div>
            <div className={classes.title}>{props.title}</div>
            <div>{props.description}</div>
            <div className={classes.btnWrapper}>
                {!going && <div id={props.id} className={classes.btn} onClick={handleGoing}>Going</div>}
                {going && <SmallButton>Going</SmallButton>}
                {!notGoing && <div id={props.id} className={classes.btn} onClick={handleNotGoing}>Not Going</div>}
                {notGoing && <SmallButton>Not Going</SmallButton>}
            </div>
        </div>
    </Card>

}

export default GroupEvent;