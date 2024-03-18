import { useContext, useState } from "react";
import { GroupsContext } from "../store/groups-context";
import { JoinedGroupContext } from "../store/joined-group-context";
import Card from "../UI/Card";
import SmallButton from "../UI/SmallButton";

import classes from './CreateGroup.module.css';

function CreateGroup(props) {

    const currUserId = localStorage.getItem("user_id");
    const currId = parseInt(currUserId);

    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');

    const grpCtx = useContext(GroupsContext);
    const jgrpCtx = useContext(JoinedGroupContext);

    console.log({title})

    function submitHandler(event) {

        // console.log("sssdsdeqfe")
        event.preventDefault();

        const date =  Date.now()

        // const created = new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short', year: '2-digit' }).format(date);

        const data = {
            id: Date.now(),
            title: title,
            creator: currId,
            description: description,
            createdat: Date.now(),
        };

        console.log(data)

        setTitle('');
        setDescription('');
    
        fetch('http://localhost:8080/group', {
            method: 'POST',
            credentials: "include",
            mode: "cors",
            body: JSON.stringify(data),
            headers: { 
                'Content-Type': 'application/json' 
            }})
            .then(resp => resp.json())
            .then(data => {
                console.log("create grp data", data);
                if (data.success) {
                    console.log("created grp resp: ", data.success);
                    grpCtx.onNewGroupCreated();
                    jgrpCtx.join(data.createdid, data.creator);
                    props.onnewgroup()
                }
            })
            .catch(err => console.log(err))     
            
    }

    return <Card className={classes.card}>
        Create Group
            <form className={classes.container} onSubmit={submitHandler}>
        <input type="text" name="title" id="title" placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} required></input>
        <textarea className={classes.content} name="description" id="description" placeholder="Description" value={description} onChange={e => setDescription(e.target.value)} ></textarea>
        <div className={classes.btn}>
            <SmallButton>Create</SmallButton> 
        </div>
        
    </form>
    </Card>

}

export default CreateGroup;