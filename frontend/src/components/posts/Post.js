import { useCallback, useEffect, useState } from 'react';
import { Link, useNavigate } from "react-router-dom";
import classes from './Post.module.css'
import AllComments from "./comments/AllComments";
import CreateComment from './comments/CreateComment';
import Avatar from '../UI/Avatar';
import Card from '../UI/Card';
import lock from '../assets/lock2.svg'
import lockCF from '../assets/lock6.svg'


function Post(props) {
    const [showComments, setShowComments] = useState(false);
    const [postPrivacy, setPostPrivacy] = useState();
    const navigate = useNavigate();

    useEffect(() => {
        if (props != undefined) {
            setPostPrivacy(props.postPrivacy)
        }
    }, [props])

    // console.log("comment for post: ", props.postNum, " comments: ", props.commentsForThisPost)
    // const onlineStatus = false;
    const postCommentUrl = "http://localhost:8080/post-comment";

    // return <div className={classes.container}>
    const showCommentsHandler = useCallback(() => {
        console.log(showComments);
        !showComments && setShowComments(true);
        showComments && setShowComments(false);
    }, []);

    const createCommentHandler = (createCommentPayloadObj) => {
        console.log("create comment for Post", createCommentPayloadObj)

        const reqOptions = {
            method: "POST",
            body: JSON.stringify(createCommentPayloadObj),
        }
        fetch(postCommentUrl, reqOptions)
            .then(resp => resp.json())
            .then(data => {
                props.onCreateCommentSuccessful(data.success); // lift it up to PostPage
            })
            .catch(err => {
                console.log(err);
            })
    };

    function handleClick(e) {
        const id = e.target.id

        console.log("id: ", id)
        navigate("/profile", { state: { id } })
    }

    // var myDate = new Date(props.createdat);
    const moment = require(`moment`);
    const myDate = moment(props.createdat)
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


    let privacy;
    if (props.privacy == 1) {
        privacy = <div> <img src={lock} alt=''></img> </div>
    } else if (props.privacy == 2) {
        privacy = <div> <img src={lockCF} alt=''></img> </div>
    }

    return <Card className={classes.container} >
        <div className={classes["author"]}>
            <Link to={`/profile/${props.authorId}`}>
                <Avatar className={classes["post-avatar"]} id={props.authorId} src={props.avatar} alt="" showStatus={false} height={"50px"} width={"50px"} />
            </Link>
            <div className={classes.column}>
                <div className={classes.row}>
                    <Link to={`/profile/${props.authorId}`}>
                        <div className={classes["details"]}>{`${props.fname} ${props.lname}`}</div>
                    </Link>
                    {privacy}
                </div>

                <div className={classes["create-at"]}>{newDate}</div>

            </div>
        </div>
        {/* <div className={classes["create-at"]}>{props.privacy}</div> //privacy */}
        <div className={classes.content}>{props.message}</div>

        {props.image && <div><img src={props.image} alt="" width={"100px"} /></div>}
        <div className={classes.comments} onClick={showCommentsHandler}>{props.commentsForThisPost.length} comments</div>
        {/* {props.commentsForThisPost.length}  */}
        {showComments &&
            <>
                <AllComments comments={props.commentsForThisPost} />
                <CreateComment pid={props.id} onCreateComment={createCommentHandler} />
            </>
        }
    </Card>
}

export default Post