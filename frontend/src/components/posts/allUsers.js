import { useContext } from 'react';
// import classes from '../group/JoinedGroup.module.css';
import classes from './allUsers.module.css'
import { useNavigate } from "react-router-dom";
import useGet from '../fetch/useGet'
import Card from '../UI/Card';
import { JoinedGroupContext } from '../store/joined-group-context';
import Avatar from '../UI/Avatar';
import SmallAvatar from '../UI/SmallAvatar';


function AllUsers() {
    const navigate = useNavigate();

    const currUserId = localStorage.getItem("user_id");

    const { error, isLoaded: isLoading, data } = useGet(`/user`)
    if (isLoading) return <div>Loading All Users...</div>
    if (error) return <div>Error: {error.message}</div>


    function handleClick(e) {
        let id = e.target.id

        if (id == "") {
            id = e.target.parentElement.parentElement.id
        }
        navigate(`/profile/${id}`)

    }

    return <>
        <div className={classes.label}>
            Users
        </div>
        <Card>
    {data && data.data && data.data.map((user) => (

                <div key={user.id} id={user.id} className={classes.container} onClick={handleClick} >
                    {/* <div className={classes.img}></div> */}
                    {/* <Avatar height={40} width={40}></Avatar> */}
                    <SmallAvatar src={user.avatar}height={35} width={35}></SmallAvatar>
                    <div>
                        <div className={classes.title}>{user.fname} {user.lname}</div>
                    </div>

                </div>
    ))}

        </Card>
    </>
}

export default AllUsers;