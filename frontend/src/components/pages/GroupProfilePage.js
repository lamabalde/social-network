import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import AllEvents from "../group/AllEvents";
import AllGroupPosts from "../group/AllGroupPosts";
import CreateEvent from "../group/CreateEvent";
import CreateGroup from "../group/CreateGroup";
import CreateGroupPost from "../group/CreateGroupPost";
import GroupEvent from "../group/GroupEvent";
import GroupProfile from "../group/GroupProfile";
import classes from './layout.module.css';
import refreshImg from "../assets/refresh.svg"
import useGet from "../fetch/useGet";

function GroupProfilePage() {
    const { state } = useLocation();
    const { id } = state;
    const [ postData, setPostData ] = useState([])
    const [refreshState, setRefreshState] = useState(false)


    useEffect(() => {
        fetch(`http://localhost:8080/group-post?groupid=${id}`)
            .then(resp => {
                return resp.json();
            })
            .then(data => {

                data.data.sort((a, b) => Date.parse(b.createdat) - Date.parse(a.createdat));
                setPostData(data.data)
            })
            .catch(
                err => console.log(err)
            );
    }, [refreshState]);

    function refresh() {
        refreshState ? setRefreshState(false) : setRefreshState(true)
    }

    function onCreatePostHandler(postData) {

        fetch('http://localhost:8080/group-post',
        {
            method: 'POST',
            credentials: "include",
            mode: "cors",
            body: JSON.stringify(postData),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(() => {
            console.log("posted")

                fetch(`http://localhost:8080/group-post?groupid=${id}`)
                .then(resp => {
                    return resp.json();
                })
                .then(data => {

                    data.data.sort((a, b) => Date.parse(b.createdat) - Date.parse(a.createdat));
                    setPostData(data.data)
                })
                .catch(
                    err => console.log(err)
                );

        })

    }
console.log(postData, "grouppostdata")
    return (
    <div className={classes.container}>
        <div className={classes.mid}>
            <GroupProfile groupid={id}></GroupProfile>
            <CreateGroupPost   groupid={id} onCreatePost={onCreatePostHandler}/>
            <div className={classes.refreshContainer}>

            <div className={classes.refresh} onClick={refresh}><img src={refreshImg} alt=""></img></div>
            </div>


            {postData &&
            <AllGroupPosts groupid={id} posts={postData}/>
            }

        </div>
        <div className={classes.right}>
        <CreateEvent groupid={id} newEvent={refresh} ></CreateEvent>
        <AllEvents groupid={id} refresh={refreshState}></AllEvents>
        </div>
    </div>

)}

export default GroupProfilePage;
