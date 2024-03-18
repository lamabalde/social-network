import { Link } from "react-router-dom";
import classes from './profilePost.module.css'
// import profile from '../assets/profile.svg';
// import AllComments from "./comments/AllComments";
// import CreateComment from './comments/CreateComment';
import Avatar from '../UI/Avatar';
import Card from '../UI/Card';
import useGet from '../fetch/useGet';
// import AllComments from './comments/AllComments';


function ProfilePosts() {
    const { error, isLoading, data } = useGet("/post");
    console.log("profile posts", data);
   
    if (isLoading) return <div>Loading...</div>
    if (error) return <div>Error: {error.message}</div>

    return <>
    {data && data.map((post) => (
    <Card className={classes.container} >
        <div className={classes["author"]}>     
            <div>
                {<Avatar className={classes["post-avatar"]} src={post.avatar} alt="" width={"50px"}/>}
            </div>
            <Link to={`/profile/${post.authorId}`}>
                <div><p className={classes["details"]}>{`${post.fname} ${post.lname} ${post.nname}`}</p></div>
            </Link>
        </div>
        <div className={classes["create-at"]}>{post.createdat.split(".")[0]}</div>
        <div className={classes.content}>{post.message}</div>
        {post.image && <div><img src={post.image} alt="" width={"100px"}/></div>}
    </Card>
        ))}
   </>
}

export default ProfilePosts