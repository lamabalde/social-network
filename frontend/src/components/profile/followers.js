import Card from "../UI/Card";
import profile from '../assets/profile.svg';
import classes from './followers.module.css'
import useGet from "../fetch/useGet";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";


function Followers({ userId }) {
    const navigate = useNavigate();

    const { error, isLoaded: isLoading, data } = useGet(`/user-follower?id=${userId}`);

    if (isLoading) return <div>Loading..foler.</div>
    if (error) return <div>Error: {error.message}</div>

    function handleClick(e) {
        const id = e.target.id

        console.log("id: ", id)
        navigate(`/profile/${id}`)

    }
    
    return <Card>
        Followers

        {data.data && data.data.map((follower) => (

            <div key={follower.id} className={classes.container}>
                <div className={classes.left}>
                    <img className={classes.img} src={profile} />
                    <div key={follower.id} id={follower.id} onClick={handleClick} className={classes.user}>{follower.fname}</div>
                </div>

                {/* <div className={classes.right}>
                </div> */}
            </div>

        ))}
    </Card>
}

export default Followers;