import { useState, useRef } from 'react'
import { Link } from "react-router-dom";
import send from '../../assets/send2.svg'



import profile from '../../assets/profile.svg'
import ImgUpload from '../../UI/ImgUpload'
import classes from './CreateComment.module.css'
import Avatar from '../../UI/Avatar'
import SmallAvatar from '../../UI/SmallAvatar';



function CreateComment(props) {
    const userId = +localStorage.getItem("user_id");
    // const first = localStorage.getItem("fname");
    // const last = localStorage.getItem("lname");
    // const nickname = localStorage.getItem("nname");
    const avatar = localStorage.getItem("avatar");

    const [uploadedCommentImg, setUploadedCommentImg] = useState("");
    const commentInput = useRef();
    // const [commentMsg, setCommentMsg] = useState("");

    function SubmitHandler(event) {
        event.preventDefault();

        const enteredContent = commentInput.current.value

        const commentData = {
            postId: props.pid,
            userId: userId, // author
            message: enteredContent,
            image: uploadedCommentImg
        };

        // console.log("create comment: ", commentData)

        props.onCreateComment(commentData);

        commentInput.current.value = "";
        setUploadedCommentImg("");
    }

    const CommentImgUploadHandler = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.addEventListener("load", () => {
            console.log("comment img", reader.result);
            setUploadedCommentImg(reader.result);
        });
    };

    return <form className={classes.inputWrapper} onSubmit={SubmitHandler}>
        <Link to={`/profile/${userId}`} className={classes["author"]}>
            <Avatar src={localStorage.getItem("avatar")}className={classes["avatar"]} id={userId}  alt="" showStatus={false} height={"35px"}  width={"35px"} />
        </Link>
        <textarea className={classes.input} placeholder="Write a comment" ref={commentInput} required />
        <div className={classes["functions"]}>
            <div className={classes["attach"]}>
                {!uploadedCommentImg && <ImgUpload name={`comment-image-${props.pid}`} id={`comment-image-${props.pid}`} accept=".jpg, .jpeg, .png, .gif" text="Attach" onChange={CommentImgUploadHandler} />}
            </div>
            <button className={classes.send}>
                <img src={send} alt='' />
            </button>
        </div>

        {uploadedCommentImg &&
            <figure className={classes["comment-img-preview"]}>
                <img src={uploadedCommentImg} height={"35px"} />
            </figure>
        }
    </form>
}

export default CreateComment;
