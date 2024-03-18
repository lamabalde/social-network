import { useContext, useEffect, useState } from 'react';
import classes from './JoinedGroup.module.css';
import { useNavigate } from "react-router-dom";
import Card from '../UI/Card';
import { JoinedGroupContext } from '../store/joined-group-context';
import groupimg from '../assets/groupimg.jpg'



function JoinedGroups({ refresh }) {
    const navigate = useNavigate();
    const jGrpCtx = useContext(JoinedGroupContext);
    const [joinedGroup, setJoinedGroup] = useState()
    const currUserId = localStorage.getItem("user_id");
    
    function handleClick(e) {
        let id = e.target.id
        // console.log("3690",e)
        if (id == "") {
            id = e.target.parentElement.parentElement.id
        }
        console.log("3690", id)
        navigate("/groupprofile", { state: { id } })
    }
    
    useEffect(() => {
        // setJoinedGroup(jGrpCtx.joinedGrps)
        if (jGrpCtx.joinedGrps && jGrpCtx.joinedGrps.length > 0){
            // alert("dssdg")
            setJoinedGroup(true)
            console.log(joinedGroup)
        }
    }, [refresh]);
    useEffect(() => {
        // setJoinedGroup(jGrpCtx.joinedGrps)
        if (jGrpCtx.joinedGrps && jGrpCtx.joinedGrps.length > 0){
            // alert("dssdg")
            setJoinedGroup(true)
            console.log(joinedGroup)
        }
    }, [jGrpCtx.joinedGrps]);
// console.log("joirndgrups",joinedGroup.length)
    return <div>
        {joinedGroup && <div className={classes.label}>
                Groups you've joined
            </div>
        }
        {joinedGroup &&  
        <Card className={classes.joinedGroupCard}>
            {/* {data.data && data.data.map((group) => ( */}
            {jGrpCtx.joinedGrps && jGrpCtx.joinedGrps.map((group) => (
                <div key={group.id} id={group.id} className={classes.container} onClick={handleClick} >
                    {/* <div className={classes.img}></div> */}
                    <div className={classes.groupimg}> <img src={groupimg} ></img> </div>
                    <div>
                        <div className={classes.title}>{group.title}</div>
                    </div>
                    {console.log("title jg", group.title)}
                </div>
            ))}
        </Card>}
    </div>
}

export default JoinedGroups;