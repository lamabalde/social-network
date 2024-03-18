import Card from "../UI/Card";
import JoinedGroup from "./JoinedGroup";
import classes from './AllJoinedGroups.module.css';
import useGet from "../fetch/useGet";

function AllJoinedGroups() {

    const currUserId = localStorage.getItem("user_id");

    const { error , isLoading, data } = useGet(`/group-member?userid=${currUserId}`)
    // for group members `/group-member?groupid=${groupId}

      if (isLoading) return <div>Loading...</div>
      if (error) return <div>Error: {error.message}</div>

    return <Card>
        <div className={classes.label}>
        Groups you've joined
        </div>
        {data.data && data.data.map((group) => (
         <JoinedGroup
        key={group.id}
        id={group.id}
        title={group.title} 
        creator={group.creator}
        description={group.description}  
        // img={group.img}
        />
        ))} 
    </Card>
}

export default AllJoinedGroups;