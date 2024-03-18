import Card from "../UI/Card";
import profile from '../assets/profile.svg';
import classes from './followers.module.css'
import useGet from "../fetch/useGet";
import { useNavigate } from "react-router-dom";


function Following({ userId }) {

    const navigate = useNavigate();

    const currUserId = localStorage.getItem("user_id");

    const { error , isLoaded: isLoading, data } = useGet(`/user-following?id=${userId}`)

    if (isLoading) return <div>Loading..foing.</div>
    if (error) return <div>Error: {error.message}</div>

    console.log("following: ", data.data);

    function handleClick(e) {
        const id = e.target.id

        console.log("id: ", id)
        navigate(`/profile/${id}`)
    }

    return <Card>
    Following
    {data.data && data.data.map((following) => (
         <div key={following.id} className={classes.wrapper}>
         <img className={classes.img} src={profile}/>
         <div key={following.id} id={following.id} onClick={handleClick} className={classes.user}>{following.fname}</div>
        </div>
     
     ))} 
    </Card>
}

export default Following;