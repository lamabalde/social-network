import { useState, useRef } from 'react'
import { Link } from "react-router-dom";
import send from '../assets/send2.svg'
import profile from '../assets/profile.svg'
import classes from './CreateGroupComment.module.css'
import Avatar from '../UI/Avatar'
import SmallAvatar from '../UI/SmallAvatar';


function CreateGroupComment(props) {
    const userId = +localStorage.getItem("user_id");
    const avatar = localStorage.getItem("avatar");

    const commentInput = useRef();

    function SubmitHandler(event) {
        event.preventDefault();

        const enteredContent = commentInput.current.value

        const commentData = {
            postId: props.pid,
            userId: userId, // author
            message: enteredContent,
            image:localStorage.getItem("avatar")
        };

        props.onCreateComment(commentData, props.pid);

        commentInput.current.value = "";
    }


    return <form className={classes.inputWrapper} onSubmit={SubmitHandler}>
        <Link to={`/profile/${userId}`} className={classes["author"]}>
            <SmallAvatar className={classes["avatar"]} src={localStorage.getItem("avatar")} alt="" height={"35px"} width={"35px"} />
        </Link>
        <textarea className={classes.input} placeholder="Write a comment" ref={commentInput} required/>
        <div className={classes["functions"]}>
            <button className={classes.send}>
                <img src={send} alt='' />
            </button>
        </div>

    </form>
}

export default CreateGroupComment;
