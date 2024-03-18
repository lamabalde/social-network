import { useNavigate } from "react-router-dom";
import { GroupsContext } from "../store/groups-context";
import { JoinedGroupContext } from "../store/joined-group-context";
import { WebSocketContext } from "../store/websocket-context";
import { UsersContext } from "../store/users-context";
import Card from "../UI/Card";
import SmallButton from "../UI/SmallButton";

import classes from './Group.module.css';
import { useContext, useEffect, useState } from "react";

function Group(props) {
    const navigate = useNavigate();

    const currUserId = localStorage.getItem("user_id");
    // console.log("curr id", currUserId);
    const [currentlyJoined, setCurrentlyJoined] = useState(false);
    const [requestedToJoin, setRequestedToJoin] = useState(false);

    const jGrpCtx = useContext(JoinedGroupContext);
    const grpCtx = useContext(GroupsContext);
    const wsCtx = useContext(WebSocketContext);
    const usersCtx = useContext(UsersContext);
    useEffect(() => {
        jGrpCtx.getJoinedGrps()
    }, [])
    useEffect(() => {
        // const storedJoinedGrps = JSON.parse(localStorage.getItem("joined-grps"));
        // console.log("stored joined-grps (group)", storedJoinedGrps);
        // check if the current group item is one of the joined Grps in the joinedGrps array
        // jGrpCtx.getJoinedGrps()


        jGrpCtx.joinedGrps && setCurrentlyJoined(jGrpCtx.joinedGrps.some(joinedGrp => joinedGrp.id === +props.grpid))
        console.log("fklglfl123", jGrpCtx.joinedGrps)
    }, [jGrpCtx.joinedGrps, props.grpid, props])

    useEffect(() => {
        jGrpCtx.requestLocalStrg()
    }, [])
    useEffect(() => {
        // const reqGroups = JSON.parse(localStorage.getItem("requestedGroups"));
        const reqGroups = jGrpCtx.requestedGroups
        console.log("stored req-grps (group)", reqGroups);
        if (reqGroups != null) {
            reqGroups.forEach(element => {
                if (element.groupid == props.grpid) {
                    setRequestedToJoin(true)
                }
            }
            )
        }
    }, [jGrpCtx.requestedGroups, props.grpid])

    useEffect(() => {
        if (wsCtx.newNotiInvitationReplyObj) {
            if (wsCtx.newNotiInvitationReplyObj.accepted) {
                setCurrentlyJoined(true);
                setRequestedToJoin(false);

                const UserJoining = usersCtx.users.find(user => user.id === wsCtx.newNotiInvitationReplyObj.sourceid);
                const JoinGroup = grpCtx.groups.find(group => group.id === wsCtx.newNotiInvitationReplyObj.groupid);
                console.log("found group to join (group accepted req)", JoinGroup);
                jGrpCtx.join(JoinGroup);

                console.log("join group id", wsCtx.newNotiInvitationReplyObj.groupid);
                console.log("cur user has joined these groups (group)", jGrpCtx.joinedGrps);
                console.log("This user join the group", UserJoining)

                joinGrpHandler(JoinGroup, UserJoining);
            } else {
                setCurrentlyJoined(false);
                setRequestedToJoin(false);
            }
        }
        wsCtx.setNewNotiJoinReplyObj(null);
    }, [wsCtx.newNotiInvitationReplyObj]);

    function reqToJoinHandler(e) {
        const grpid = e.target.id;
        console.log("grpid", e.target.id);
        jGrpCtx.requestToJoin(+grpid);
        jGrpCtx.requestLocalStrg();
        setRequestedToJoin(true);
        setCurrentlyJoined(false);

        const data = {
            id: Date.now(),
            userid: parseInt(currUserId),
            groupid: parseInt(grpid),
            status: 0,
            createdat: Date.now(),
        };

        console.log(data)

        fetch('http://localhost:8080/group-request',
            {
                method: 'POST',
                credentials: "include",
                mode: "cors",
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(() => {
                // navigate.replace('/??')
                console.log("group request posted")
            })
    }

    const joinGrpHandler = (grp, user) => {
        console.log("user joining group (group)", grp);
        // jGrpCtx.storeGroupMember(grp, user);
        setCurrentlyJoined(true);
        setRequestedToJoin(false);
    };


    function handleClick(e) {
        const id = e.target.id
        console.log(id, "1234536546")
        navigate("/groupprofile", { state: { id } })

    }
    return <Card>
        <div className={classes.container}>
            <div className={classes.wrapper}>
                <div className={classes.img}></div>
                <div>
                    <div className={classes.title}>{props.title}</div>
                    <div className={classes.desc}>{props.description}</div>
                </div>

            </div>
            <div className={classes.btn}>
                {!currentlyJoined && !requestedToJoin && <div className={classes.smallbtn} id={props.grpid} onClick={reqToJoinHandler}>Join</div>}
                {!currentlyJoined && requestedToJoin && <div className={classes.smallbtn} id={props.grpid}>Requested</div>}
                {currentlyJoined && !requestedToJoin && <div className={classes.smallbtn} id={props.grpid} onClick={handleClick} >Joined</div>}
                {/* {currentlyJoined && requestedToJoin && <div className={classes.smallbtn} id={props.grpid}>Requested & Joined</div>} */}
            </div>
        </div>
    </Card>
}

export default Group;