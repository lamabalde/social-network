import classes from './GroupPost.module.css'
import profile from '../assets/profile.svg';
import Card from '../UI/Card';
import { useState } from 'react';
import AllComments from '../posts/comments/AllComments';
import CreateGroupComment from './CreateGroupComment';
import { Link, useNavigate } from 'react-router-dom';
import Avatar from '../UI/Avatar';
import SmallAvatar from '../UI/SmallAvatar';

function GroupPost(props) {

    const [commentData, setCommentData] = useState([]);
    const [showComments, setShowComments] = useState(false);

    const showCommentsHandler = (e) => {
        console.log("props id target: ", e.target.id)

        !showComments && setShowComments(true);
        showComments && setShowComments(false);

        fetch(`http://localhost:8080/group-post-comment?id=${e.target.id}`)
            .then(resp => resp.json())
            .then(data => {
                setCommentData(data.data);
                console.log("comments data", data)
            })
            .catch(
                err => console.log(err)
            );
    };

    const createCommentHandler = (createCommentPayloadObj) => {

        const reqOptions = {
            method: "POST",
            body: JSON.stringify(createCommentPayloadObj),
        }
        fetch("http://localhost:8080/group-post-comment", reqOptions)
            .then(resp => resp.json())
            .then(data => {
                createCommentSuccessHandler(data.success);
            })
            .catch(err => {
                console.log(err);
            })
    };

    const createCommentSuccessHandler = (createCommentSuccessful) => {
        // fetch comment
        console.log({createCommentSuccessful})
        console.log("props id: ", props.id)
        if (createCommentSuccessful) {
        fetch(`http://localhost:8080/group-post-comment?id=${props.id}`)
            .then(resp => resp.json())
            .then(data => {
                console.log("post page raw comment data: ", data.data)

                setCommentData(data.data);
            })
            .catch(
                err => console.log(err)
            );
        }
    };

    const created = new Intl.DateTimeFormat('en-GB', {
        day: 'numeric',
        month: 'short',
        year: '2-digit',
        hour: 'numeric',
        minute: 'numeric'
      }).format(new Date(props.createdat));
       console.log("created ", created)
console.log(props, "porpos.avatar")
    return <Card className={classes.container} >
        <div className={classes.user}>
            {/* <img src={profile} alt='' /> */}
            <SmallAvatar src={props.avatar}width={40}  height={40}></SmallAvatar>
            <div>
                <Link to={`/profile/${props.author}`} className={classes.username}>{props.fname} {props.lname}</Link>
                <div className={classes.date}>{created}</div>
            </div>

        </div>
        <div className={classes.content}>{props.message}</div>
        <div className={classes.comments} id={props.id} onClick={showCommentsHandler}> comments</div>
        {/* {commentData.length} */}
        {showComments &&
            <>
                {commentData &&
                  <AllComments comments={commentData} />
                   }

                <CreateGroupComment pid={props.id} onCreateComment={createCommentHandler} />
            </>
        }
    </Card>

}

export default GroupPost