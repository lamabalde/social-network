import classes from './GroupProfile.module.css';
import SmallButton from "../UI/SmallButton";
import GreyButton from "../UI/GreyButton";
import Card from "../UI/Card";
import useGet from '../fetch/useGet';
import { useLocation, useNavigate } from 'react-router-dom';
import Modal from './modal';
import { useState, useContext, useEffect } from 'react';
import { GroupsContext } from "../store/groups-context";
import { JoinedGroupContext } from "../store/joined-group-context";
import { WebSocketContext } from "../store/websocket-context";
import { UsersContext } from "../store/users-context";
import groupimg from '../assets/groupimg.jpg'

function GroupProfile({ groupid }) {

    const navigate = useNavigate();

    const [currentlyJoined, setCurrentlyJoined] = useState(false);
    const [invitedToJoin, setInvitedToJoin] = useState(false);
    const [groupMembers, setgroupmembers] = useState(null);

    const jGrpCtx = useContext(JoinedGroupContext);
    const grpCtx = useContext(GroupsContext);
    const wsCtx = useContext(WebSocketContext);
    const usersCtx = useContext(UsersContext);

    useEffect(() => {
        console.log(wsCtx.newNotiInvitationReplyObj);
        if (wsCtx.newNotiInvitationReplyObj) {
            if (wsCtx.newNotiInvitationReplyObj.accepted) {
                setCurrentlyJoined(true);
                setInvitedToJoin(false);

                const UserJoining = usersCtx.users.find(user => user.id === wsCtx.newNotiInvitationReplyObj.sourceid);
                const JoinGroup = grpCtx.groups.find(group => group.id === wsCtx.newNotiInvitationReplyObj.groupid);
                console.log("found group to join (invited) (group accepted req)", JoinGroup);
                jGrpCtx.join(JoinGroup);

                console.log("join group id (invited)", wsCtx.newNotiInvitationReplyObj.groupid);
                console.log("this user (invited) to join the group", UserJoining)

                joinGrpHandler(JoinGroup, UserJoining);
            } else {
                setCurrentlyJoined(false);
                setInvitedToJoin(false);
            }
        }
        wsCtx.setNewNotiInvitationReplyObj(null);
    }, [wsCtx.newNotiInvitationReplyObj]);

    const { error, isLoaded: isLoading, data } = useGet(`/group?id=${groupid}`)
    const [open, setOpen] = useState(false)


    const currUserId = localStorage.getItem("user_id");

    if (isLoading) return <div>Loading GP...</div>
    if (error) return <div>Error: {error.message}</div>

    function handleClick(e) {
        const id = e.target.id;
        const currUserId = localStorage.getItem("user_id");

        setOpen(true)

        navigate("/groupprofile", { state: { id } })

    }

    const joinGrpHandler = (grp, user) => {
        console.log("user joining group (invite)", grp);
        jGrpCtx.storeGroupMember(grp, user);
        setCurrentlyJoined(true);
        setInvitedToJoin(false);
    };

    const invitedHandler = (invited) => {
        if (invited) {
            setInvitedToJoin(true);
            setCurrentlyJoined(false);
        } else {
            setInvitedToJoin(false)
        }
    };

    return <Card className={classes.container}>
        {data && data.data && data.data.map((group) => (
            <div className={classes.groupContainer} key={group.id} id={group.id}>
                <div className={classes.img}> <img src={groupimg} ></img> </div>

                <div className={classes.wrapper}>
                    <div className={classes.row}>
                        <div className={classes.groupname}>{group.title}</div>


                        <div className={classes.btnWrapper}>
                            <div id={group.id} className={classes.btn} onClick={handleClick}>+ Invite</div>
                            <GreyButton>Message</GreyButton>
                        </div>
                    </div>

                    <div className={classes.description}>{group.description}</div>
                    {/* <div className={classes.members}>Members</div> */}
                </div>
                <Modal
                    open={open}
                    onClose={() => setOpen(false)}
                    onInvite={invitedHandler}
                    currentlyJoined={currentlyJoined}
                    invitedToJoin={invitedToJoin}
                ></Modal>
            </div>
        ))}
    </Card>
}

export default GroupProfile;