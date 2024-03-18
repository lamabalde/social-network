import { useRef, useState } from "react";
import { Link } from "react-router-dom";
import Card from "../UI/Card";
import CreatePostTextarea from "../UI/CreatePostTextarea";
import SmallButton from "../UI/SmallButton";
import FormPostSelect from "../UI/FormPostSelect";
import ImgUpload from "../UI/ImgUpload";
import Avatar from "../UI/Avatar";
import classes from './CreatePost.module.css';
import img1 from '../assets/img.svg'

function CreatePost(props) {
    const userId = +localStorage.getItem("user_id");
    const first = localStorage.getItem("fname");
    const last = localStorage.getItem("lname");
    const nickname = localStorage.getItem("nname");
    const avatar = localStorage.getItem("avatar");

    const [uploadedImg, setUploadedImg] = useState("");
   
    const contentInput = useRef();
    const privacyInputRef = useRef();

    function SubmitHandler(event) {
        event.preventDefault();

        const enteredContent = contentInput.current.value
        const chosenPrivacy = +privacyInputRef.current.value;

        const postData = {
            author: userId,
            message: enteredContent,
            image: uploadedImg,
            privacy: chosenPrivacy
        };

        console.log("create post data", postData)

        props.onCreatePost(postData)

        contentInput.current.value = "";
        privacyInputRef.current.value = 0;
        setUploadedImg("");
    }
    const imgUploadHandler = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.addEventListener("load", () => {
            console.log(reader.result);
            setUploadedImg(reader.result);
        })
    };

    const privacyOptions = [
        {value: 0, text:  "Public Post"},
        {value: 1, text: `Private Post`},
        {value: 2, text: "Close Friends"}
    ];

    return <form className={classes.container} onSubmit={SubmitHandler}>

        <Card>
            <div className={classes["author"]}>
            <Link to={`/profile/${userId}`} >
                <Avatar className={classes["avatar"]} id={userId} src={avatar} showStatus={false} height={"40px"}  alt="" width={"40px"}/>
            </Link> 
            <Link to={`/profile/${userId}`} className={classes["author"]}>
                <div className={classes["details"]}>{`${first} ${last}`}</div>
            </Link>
            </div>
            <div className={classes["content-container"]}>
                <div>
                    <CreatePostTextarea className={classes.content} placeholder="What's on your mind?" reference={contentInput} rows="3"/>
                </div>
                <div>
                <figure>
                    {uploadedImg && <img src={uploadedImg} className={classes["img-preview"]} width={"80px"}/>}
                </figure>
                    <ImgUpload className={classes["attach"]} name="image" id="image" accept=".jpg, .jpeg, .png, .gif" text="Attach" onChange={imgUploadHandler}/>
                </div>
                <div>
                    <FormPostSelect options={privacyOptions} className={classes["privacy"]} reference={privacyInputRef}/>
                </div>
            </div>
        
            <div className={classes.btn}>
            
            <SmallButton>Post</SmallButton>
         </div>
        </Card>
      
         
    </form>
}

export default CreatePost;