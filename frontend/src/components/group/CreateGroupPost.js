import { useState } from "react";
import Card from "../UI/Card";
import FormTextarea from "../UI/FormTextarea";
import SmallButton from "../UI/SmallButton";
import classes from './CreateGroupPost.module.css';
import profile from '../assets/profile.svg';
import Avatar from "../UI/Avatar";
import SmallAvatar from "../UI/SmallAvatar";


function  CreateGroupPost(props) {

    const currUserId = localStorage.getItem("user_id");
    const [ message, setMessage ] = useState('')

    function SubmitHandler(event) {
        event.preventDefault();
        setMessage('');
        const data = {
            id: 0,
            author: parseInt(currUserId),
            groupid: parseInt(props.groupid),
            message: message,
            image: localStorage.getItem("avatar"),
            createdat: Date.now(),
        };

        props.onCreatePost(data)


        // fetch('http://localhost:8080/group-post',
        // {
        //     method: 'POST',
        //     credentials: "include",
        //     mode: "cors",
        //     body: JSON.stringify(data),
        //     headers: {
        //         'Content-Type': 'application/json'
        //     }
        // }).then(() => {
        //     // navigate.replace('/??')
        //     console.log("posted")

        //     if (data.success) {
        //         // render all posts
        //         fetch(`http://localhost:8080/group-post?groupid=${props.groupid}`)
        //         .then(resp => {
        //             return resp.json();
        //         })
        //         .then(data => {
        //             // console.log("post data: ", data)
        //             // setPostData(data) // no need to sort
        //             // console.log("parsed 0", Date.parse(data[0].createdat));
        //             // console.log("parsed 1", Date.parse(data[1].createdat));
        //             data.sort((a, b) => Date.parse(b.createdat) - Date.parse(a.createdat));
        //             console.log("sorted post data: ", data);
        //             props.onCreatePost(data)
        //         })
        //         .catch(
        //             err => console.log(err)
        //         );
        //     }


        // })
    }

return <form className={classes.container} onSubmit={SubmitHandler}>
        <Card>
            <div className={classes.row}>
                {/* <img src={profile} alt='' /> */}
                <SmallAvatar src={localStorage.getItem("avatar")}height={40} width={40}/>
                <textarea className={classes.content} placeholder="What's on your mind?" rows="1" value={message} onChange={e => setMessage(e.target.value)} required/>
            </div>

        <div className={classes.btn}>
            <SmallButton>Post</SmallButton>
        </div>
        </Card>
    </form>
}

export default CreateGroupPost;
