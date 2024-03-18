import { useContext,useState } from "react";
import Card from "../UI/Card";
import SmallButton from "../UI/SmallButton";
import { WebSocketContext } from "../store/websocket-context";
import classes from './CreateEvent.module.css';
function CreateEvent( {groupid, newEvent} ) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [date, setDate] = useState('');
    const wsCtx = useContext(WebSocketContext);

    function SubmitHandler(event) {
        event.preventDefault();

        const currUserId = localStorage.getItem("user_id");

        const datenow =  Date.now().toString();

        const data = {
            id: 0,
            groupid: parseInt(groupid),
            author: parseInt(currUserId),
            title: title,
            description: description,
            createdat: datenow,
            date: date
        };
  console.log({data})
        const followPayloadObj = {};
        followPayloadObj.label = "noti";
        followPayloadObj.id = Date.now();
        followPayloadObj.type = "event-notif";
        followPayloadObj.sourceid = parseInt(currUserId);
        followPayloadObj.groupid = parseInt(groupid);
        followPayloadObj.targetid = 987;
        followPayloadObj.createdat = datenow;
        console.log("CREATED AT: ",followPayloadObj.createdat)
        if (wsCtx.websocket !== null) wsCtx.websocket.send(JSON.stringify(followPayloadObj));
        setTitle('');
        setDescription('');
        setDate('');
    
        fetch('http://localhost:8080/group-event', 
        {
            
            method: 'POST',
            credentials: "include",
            mode: 'cors',
            body: JSON.stringify(data),
            headers: { 
                'Content-Type': 'application/json' 
            }
        }).then(() => {
            newEvent()
            console.log("event posted")
        })
    }

    return <Card className={classes.card}>
        Create Event
            <form className={classes.container} onSubmit={SubmitHandler}>
        <input type="text" name="title" id="title" placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} required></input>
        <textarea className={classes.content} name="description" id="description" placeholder="Description" value={description} onChange={e => setDescription(e.target.value)} required></textarea>
        <input type="datetime-local" name="date" id="date" value={date} onChange={e => setDate(e.target.value)} required></input>
        <div className={classes.btn}>
            <SmallButton>Create</SmallButton>
        </div>
    </form>
    </Card>
}

export default CreateEvent;
