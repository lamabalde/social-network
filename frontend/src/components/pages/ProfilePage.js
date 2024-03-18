import { useContext, useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import useGet from "../fetch/useGet";
import AllPosts from "../posts/AllPosts"
import CreatePost from "../posts/CreatePost";
import Followers from "../profile/followers";
import Following from "../profile/following";
import Profile from "../profile/Profile";
import ProfilePosts from "../profile/profilePost";
import FollowRequest from "../requests/FollowRequest";
import { FollowingContext } from "../store/following-context";
// import classes from './ProfilePage.module.css';
import { useParams } from "react-router-dom";
import classes from './layout.module.css';

function ProfilePage() {
    const [commentData, setCommentData] = useState([]);

    const sessionUrl = "http://localhost:8080/session";
    // const { state } = useLocation();
    // const { id } = state;
    const params = useParams();
    const id = params.userId;

    // get comments
    useEffect(() => {
        fetch("http://localhost:8080/post-comment")
            .then(resp => resp.json())
            .then(data => {
                // console.log("post page raw comment data: ", data)
                // setCommentData(data);
                data.sort((a, b) => Date.parse(a.createdat) - Date.parse(b.createdat)); // ascending order
                // console.log("post page sorted comment data: ", data)
                setCommentData(data);
            })
            .catch(
                err => console.log(err)
            );
    }, []);

    const { error, isLoaded: isLoading, data } = useGet(`/post?id=${localStorage.getItem("user_id")}`)

    if (isLoading) return <div>Loading...Profile Page</div>
    if (error) return <div>Error: {error.message}</div>

    let postData = []
    if (data && data.data && data.data.length > 0){
        postData = data.filter(x => x.author === +id)}


    // let curFollowing;
    // const [curFollowing, setCurFollowing] = useState(false);
    // const checkCurFollowing = () => {
    //     const storedFollowing = JSON.parse(localStorage.getItem("following"));
    //     console.log("stored following (profile)", storedFollowing);
    //     if (followingCtx.following) setCurFollowing(followingCtx.following.some(followingUser => followingUser.id === +id))
    // };

    // useEffect(() => checkCurFollowing(), [followingCtx.following]);


    return <div className={classes.container}>
        <div className={classes.mid}>
            <Profile userId={id} ></Profile>
            <AllPosts userId={id} posts={postData} comments={commentData}></AllPosts>
        </div>
        {/* <div className={classes.right}>
            <Followers userId={id}></Followers>
            <Following userId={id}></Following>
        </div> */}
    </div>
}

export default ProfilePage;