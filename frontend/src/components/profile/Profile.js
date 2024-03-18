import { useEffect, useState, useContext } from "react";
import useGet from "../fetch/useGet";
import Card from "../UI/Card";
import GreyButton from "../UI/GreyButton";
import ToggleSwitch from "../UI/ToggleSwitch";
import { FollowingContext } from "../store/following-context";
import { UsersContext } from "../store/users-context";
import { WebSocketContext } from "../store/websocket-context";
import classes from './Profile.module.css';
import FollowerModal from "./followerModal";
import FollowingModal from "./followingModal";
import lock from '../assets/lock7.svg'
import Avatar from "../UI/Avatar";
function Profile({ userId }) {

    const [followerOpen, setFollowerOpen] = useState(false)
    const [followingOpen, setFollowingOpen] = useState(false)
    const [followerData, setFollowerData] = useState([])
    const [followingData, setFollowingData] = useState([])
    const [isFollower, setIsFollower] = useState(false)
    // get stored publicity from localStorage
    // let selfPublicStatus;
    // selfPublicNum === 0 ? selfPublicStatus = false : selfPublicStatus = true;
    let statusofcuruser;
    // self
    const [publicity, setPublicity] = useState(true); // 1 true is public, 0 false is private
    const selfPublicNum = +localStorage.getItem("public");
    const [pubCheck, setPubCheck] = useState(false)

    useEffect(() => {
        const targetUser = usersCtx.users.find(user => user.id === +userId);
        console.log("checkingTargetUser", targetUser);
        if (targetUser != undefined) {
            if (targetUser.public != 0) {
                setPubCheck(true)
            }

        }
    }, [userId]);
    console.log("stored publicity (profile)", selfPublicNum);
    useEffect(() => {
        selfPublicNum ? setPublicity(true) : setPublicity(false);
    }, [selfPublicNum]);

    // friend
    const followingCtx = useContext(FollowingContext);
    const usersCtx = useContext(UsersContext);
    const wsCtx = useContext(WebSocketContext);

    const currUserId = localStorage.getItem("user_id");

    // const wsCtx = useContext(WebSocketContext);
    // console.log("ws in profile: ",wsCtx.websocket);
    // wsCtx.websocket.onopen = function (){alert("ws is open")}
    // let params = useParams();
    // let userId = params.productId
    // console.log("params", userId)
    const [currentlyFollowing, setCurrentlyFollowing] = useState(false);
    const [requestedToFollow, setRequestedToFollow] = useState(false);
    const [isCloseFriend, setCloseFriend] = useState(false);

    useEffect(() => {
        const storedFollowing = JSON.parse(localStorage.getItem("following"));
        console.log("stored following (profile)", storedFollowing);
        // check if the current profile is one of the users in the following array
        followingCtx.following && setCurrentlyFollowing(followingCtx.following.some(followingUser => followingUser.id === +userId))
    }, [followingCtx.following, userId]);

    useEffect(() => {
        followingData && setIsFollower(followingData.some(follower => follower.id == currUserId))
    }, [followingData])

    console.log("followingData", followingData)
    // store follower in db
    const storeFollow = (targetId) => {
        console.log("targetid (storeFollow)", targetId)

        const data = {
            // id: 0,
            sourceid: parseInt(currUserId),
            targetid: parseInt(targetId),
            status: 1
        };

        const reqOptions = {
            method: "POST",
            credentials: "include",
            mode: "cors",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        };

        fetch('http://localhost:8080/user-follower', reqOptions)
            .then(resp => resp.json())
            .then(data => {
                console.log("user follower data", data);
                if (data.success) {
                    console.log("followrequest")
                }
            })
            .catch(err => {
                console.log(err);
            })
    };

    useEffect(() => {
        if (wsCtx.newNotiFollowReplyObj) {
            if (wsCtx.newNotiFollowReplyObj.accepted) {
                setCurrentlyFollowing(true);
                setRequestedToFollow(false);

                const followUser = usersCtx.users.find(user => user.id === wsCtx.newNotiFollowReplyObj.sourceid);
                console.log("found user frd (profile accepted req)", followUser);
                followingCtx.follow(followUser);

                console.log("follow user id", wsCtx.newNotiFollowReplyObj.sourceid);
                console.log("cur user is following (profile)", followingCtx.following);
                const targetId = wsCtx.newNotiFollowReplyObj.targetid;
                console.log("targetid", targetId)
                console.log("current user", currUserId)

                // storeFollow(targetId);
            } else {
                setCurrentlyFollowing(false);
                setRequestedToFollow(false);
            }
        }
        wsCtx.setNewNotiFollowReplyObj(null);
    }, [wsCtx.newNotiFollowReplyObj]);

    const followHandler = (e) => {
        const followUser = usersCtx.users.find(user => user.id === +e.target.id);
        console.log("found user frd (profile)", followUser);
        console.log("frd publicity", publicity);
        if (followUser) {
            if (followUser.public) {
                console.log(" user frd (public)");
                followingCtx.follow(followUser);
                setCurrentlyFollowing(true);
                storeFollow(e.target.id);
            } else if (!followUser.public) { //if frd private
                console.log(" user frd (private)");
                followingCtx.requestToFollow(followUser);
                setRequestedToFollow(true);
            }
        } else {
            console.log("user frd not found (profile)");
        }
    };

    const unfollowHandler = (e) => {
        console.log("unfollow userid", e.target.id);
        const unfollowUser = usersCtx.users.find(user => user.id === +e.target.id);
        console.log("unfollow user", unfollowUser);
        unfollowUser && followingCtx.unfollow(unfollowUser);
        setCurrentlyFollowing(false);

        // delete from db
        deleteFollow(e.target.id);
    };

    const setPublicityHandler = (e) => {

        if (e.target.checked) {
            setPublicity(false); // private
        } else {
            setPublicity(true);
            // usersCtx.onPrivacyChange(currUserId, 1);
        }

        let publicityNum;
        if (e.target.checked) {
            publicityNum = 0
        } else {
            publicityNum = 1;
        }
        console.log({ publicityNum })
        localStorage.setItem("public", publicityNum);

        // post to store publicity to db
        const data = {
            // Define the data to send in the request body
            targetid: parseInt(userId),
            public: publicityNum,
        };

        fetch('http://localhost:8080/privacy',
            {
                method: 'POST',
                credentials: "include",
                mode: 'cors',
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(() => {
                console.log("privacy changed")
            })
    };


    useEffect(() => {
        fetch(`http://localhost:8080/user-follow-status?tid=${userId}&sid=${currUserId}`)
            .then(response => response.text())
            .then(data => {
                console.log("cf: ", data)
                if (data == "trueclosefriend" || data == "true") {
                    setRequestedToFollow(true)
                } else {
                    setRequestedToFollow(false)
                }
                if (data == "trueclosefriend" || data == "falseclosefriend") {
                    setCloseFriend(true)
                } else {
                    setCloseFriend(false)
                }

            }).catch(error => {
                console.log({ error })
            });

    }, [userId]);


    //Get followers
    useEffect(() => {
        fetch(`http://localhost:8080/user-follower?id=${userId}`)
            .then(resp => resp.json())
            .then(data => {
                setFollowerData(data.data)
            })
            .catch(err => {
                console.log(err);
            })
    }, [userId]);

    //Get following
    useEffect(() => {
        fetch(`http://localhost:8080/user-following?id=${userId}`)
            .then(resp => resp.json())
            .then(data => {
                setFollowingData(data.data)
            })
            .catch(err => {
                console.log(err);
            })
    }, [userId]);

    const { error, isLoaded: isLoading, data } = useGet(`/user?id=${userId}`)
    if (data && data.data && data.data.length > 0) {
        if (data.data[0].public === 0) {
            localStorage.setItem('isChecked', true);
        } else {
            localStorage.setItem('isChecked', false);
        }
    }
    //  console.log("user data (profile)", data.data)
    if (isLoading) return <div>Loading Profile...</div>
    if (error) return <div>Error: {error.message}</div>

    // delete follower from db
    const deleteFollow = (targetId) => {
        console.log("targetid (deleteFollow)", targetId)

        const data = {
            // id: 0,
            sourceid: parseInt(currUserId),
            targetid: parseInt(targetId),
            status: 1
        };

        const reqOptions = {
            method: "DELETE",
            credentials: "include",
            mode: "cors",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        };

        fetch('http://localhost:8080/user-follower', reqOptions)
            .then(resp => resp.json())
            .then(data => {
                console.log(data);
                if (data.success) {
                    console.log("followrequest")
                }
            })
            .catch(err => {
                console.log(err);
            })
    };

    function closeFriendHandler(e) {
        if (isCloseFriend) {
            setCloseFriend(false)

        } else {
            setCloseFriend(true)
        }
        console.log("closefriend event: ", e)

        const data = {
            sourceid: parseInt(e.target.id),
            targetid: parseInt(currUserId),
        };
        console.log("closefriend event: ", data)
        const cfOptions = {
            method: "POST",
            credentials: "include",
            mode: "cors",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        };

        fetch('http://localhost:8080/close-friend', cfOptions)
            .then(resp => resp.json())
            .then(data => {
                console.log("closefriend event: ", data);
                if (data.success) {
                    console.log("closefriendchanges")
                }
            })
            .catch(err => {
                console.log("closefriend event: ", err);
            })

    }

    let followButton;
    let messageButton;
    let closeFriend;
    let closeFriendText;

    if (currUserId !== userId) {
        if (currentlyFollowing) {
            followButton = <div className={classes.followbtn} id={userId} onClick={unfollowHandler}>Unfollow</div>
            console.log("currentlyFollowing", currentlyFollowing);
        } else if (requestedToFollow) {
            followButton = <div className={classes.followbtn} id={userId}>Requested</div>
        } else {
            followButton = <div className={classes.followbtn} id={userId} onClick={followHandler}>+ Follow</div>
        }
        messageButton = <GreyButton>Message</GreyButton>
        closeFriend = <input className={classes.checkbox} type="checkbox" id={userId} checked={isCloseFriend} onClick={closeFriendHandler} />
        closeFriendText = <div>close friend</div>

    }

    function handleFollowerClick() {
        setFollowerOpen(true)
    }

    function handleFollowingClick() {
        setFollowingOpen(true)
    }

    return <div className={classes.container}>
        <div className={classes.private}>
            {currUserId === userId && data &&
                <ToggleSwitch
                    label={"Private"}
                    value={"Private"}
                    onClick={setPublicityHandler}
                ></ToggleSwitch>
            }
        </div>
        <Card>
            <div className={classes.wrapper}>
                {/* <div className={classes.img}></div> */}
                {/* <Avatar height={100}></Avatar> */}
                {data && 
                <>
                    <Avatar src={data.data[0].avatar} showStatus={false} height={100} width={100}/>

                    <div className={classes.column}>
                        <div className={classes.row}>
                            <div className={classes.row}>
                            <div className={classes.name}>{data.data[0].fname} {data.data[0].lname}</div>
                            {!pubCheck &&
                                <div><img src={lock} alt=''></img> </div>
                            }
                            </div>
                            <div className={classes.btn}>
                                {followButton}
                                {messageButton}
                            </div>
                        </div>

                        <div className={classes.username}>{data.data[0].nname}</div>
                        <div className={classes.followers}>
                            <div className={classes.follow} onClick={handleFollowerClick}><span className={classes.count}>{followerData && followerData.length}{!followerData && 0}</span> followers</div>
                            <div className={classes.follow} onClick={handleFollowingClick}><span className={classes.count}>{followingData && followingData.length}{!followingData && 0}</span> following</div>
                        </div>
                        <div>{data.data[0].about}</div>
                        <div className={classes.closeFriend}>
                            {isFollower &&
                                closeFriendText
                            }
                            {isFollower &&
                                closeFriend}
                        </div>
                    </div>
                </>}
            </div>
            {followerOpen && followerData &&
                <FollowerModal
                    onClose={() => setFollowerOpen(false)}
                    followers={followerData}
                ></FollowerModal>
            }
            {followingOpen && followingData &&
                <FollowingModal
                    onClose={() => setFollowingOpen(false)}
                    following={followingData}
                ></FollowingModal>
            }

        </Card>
    </div>
}


export default Profile;
